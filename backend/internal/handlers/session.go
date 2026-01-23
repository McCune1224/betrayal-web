package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

import "betrayal-web/internal/db/sqlc"

type SessionHandler struct {
	Queries *sqlc.Queries
}

type CreateSessionRequest struct {
	PlayerName string `json:"player_name"`
	RoomCode   string `json:"room_code"`
	IsHost     bool   `json:"is_host"`
	RoleID     int    `json:"role_id"`
	SessionID  string `json:"session_id,omitempty"`
	PlayerID   string `json:"player_id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

type CreateSessionResponse struct {
	SessionID  string    `json:"session_id"`
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	RoomCode   string    `json:"room_code"`
	IsHost     bool      `json:"is_host"`
	RoleID     int       `json:"role_id"`
	IsAlive    bool      `json:"is_alive"`
	CreatedAt  time.Time `json:"created_at"`
}

type RestoreSessionResponse struct {
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	RoomCode   string    `json:"room_code"`
	IsHost     bool      `json:"is_host"`
	RoleID     int       `json:"role_id"`
	IsAlive    bool      `json:"is_alive"`
	CreatedAt  time.Time `json:"created_at"`
}

func (sh *SessionHandler) CreateSession(c echo.Context) error {
	var req CreateSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Use provided SessionID if available, else generate
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	playerID := req.PlayerID
	if playerID == "" {
		playerID = uuid.New().String()
	}
	var createdAt time.Time
	if req.CreatedAt != "" {
		ct, err := time.Parse(time.RFC3339Nano, req.CreatedAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid created_at format"})
		}
		createdAt = ct
	} else {
		createdAt = time.Now().UTC()
	}

	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid uuid"})
	}
	params := sqlc.CreateSessionParams{
		ID:         sid,
		PlayerID:   playerID,
		PlayerName: req.PlayerName,
		RoomCode:   req.RoomCode,
		IsHost:     req.IsHost,
		RoleID:     sql.NullInt32{Int32: int32(req.RoleID), Valid: true},
		IsAlive:    true,
		CreatedAt:  sql.NullTime{Time: createdAt, Valid: true},
	}
	if err := sh.Queries.CreateSession(c.Request().Context(), params); err != nil {
		fmt.Printf("CreateSession error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  createdAt.Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   os.Getenv("NODE_ENV") == "production", // toggle by environment
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, CreateSessionResponse{
		SessionID:  sessionID,
		PlayerID:   playerID,
		PlayerName: req.PlayerName,
		RoomCode:   req.RoomCode,
		IsHost:     req.IsHost,
		RoleID:     req.RoleID,
		IsAlive:    true,
		CreatedAt:  createdAt,
	})
}

func (sh *SessionHandler) RestoreSession(c echo.Context) error {
	cookie, err := c.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no session"})
	}
	sid, err := uuid.Parse(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid session id"})
	}
	sess, err := sh.Queries.GetSessionByID(c.Request().Context(), sid)
	if err != nil {
		fmt.Printf("RestoreSession error: %v\n", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid session"})
	}
	resp := RestoreSessionResponse{
		PlayerID:   sess.PlayerID,
		PlayerName: sess.PlayerName,
		RoomCode:   sess.RoomCode,
		IsHost:     sess.IsHost,
		RoleID:     int(sess.RoleID.Int32),
		IsAlive:    sess.IsAlive,
		CreatedAt:  sess.CreatedAt.Time,
	}
	return c.JSON(http.StatusOK, resp)

}

func (sh *SessionHandler) DeleteSession(c echo.Context) error {
	cookie, err := c.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no session"})
	}
	sid, err := uuid.Parse(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid session id"})
	}
	if err := sh.Queries.DeleteSessionByID(c.Request().Context(), sid); err != nil {
		fmt.Printf("DeleteSession error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	expired := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   os.Getenv("NODE_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(expired)
	return c.JSON(http.StatusOK, map[string]string{"status": "session deleted"})
}
