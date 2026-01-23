package handlers

import (
	"betrayal-web/internal/db/sqlc"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setupEcho() *echo.Echo {
	e := echo.New()
	return e
}

func TestSessionHandler_CreateRestoreDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := sqlc.New(db)
	sh := &SessionHandler{Queries: queries}
	playerName := "Tester"
	roomCode := "ABCDE123"
	isHost := true
	roleID := 5
	playerID := "f4c6b5aa-aedd-4591-aaab-12aab1234567"
	sessionID := "11111111-1111-1111-1111-111111111111"
	sid, _ := uuid.Parse(sessionID)
	createdAt := time.Now().UTC()

	fmt.Printf("CreateSession args:\n sid: %T %v\n playerID: %T %v\n playerName: %T %v\n roomCode: %T %v\n isHost: %T %v\n roleID: %T %v\n isAlive: %T %v\n createdAtAny: %T %v\n",
		sid, sid, playerID, playerID, playerName, playerName, roomCode, roomCode, isHost, isHost, sql.NullInt32{Int32: int32(roleID), Valid: true}, sql.NullInt32{Int32: int32(roleID), Valid: true}, true, true, sql.NullTime{Time: createdAt, Valid: true}, sql.NullTime{Time: createdAt, Valid: true})
	mock.ExpectExec("INSERT INTO sessions").
		WithArgs(sid, playerID, playerName, roomCode, isHost, sql.NullInt32{Int32: int32(roleID), Valid: true}, true, sql.NullTime{Time: createdAt, Valid: true}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create
	e := setupEcho()
	req := httptest.NewRequest(http.MethodPost, "/api/session", strings.NewReader(
		fmt.Sprintf(`{"player_name": "%s", "room_code": "%s", "is_host": %t, "role_id": %d, "session_id": "%s", "player_id": "%s", "created_at": "%s"}`,
			playerName, roomCode, isHost, roleID, sessionID, playerID, createdAt.Format(time.RFC3339Nano))))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/session")
	_ = sh.CreateSession(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Expect select for restore
	fmt.Printf("RestoreSession args: sid: %T %v\n", sid, sid)
	mock.ExpectQuery(`SELECT id, player_id, player_name, room_code, is_host, role_id, is_alive, created_at, last_seen FROM sessions WHERE id = \$1`).WithArgs(sid).WillReturnRows(sqlmock.NewRows([]string{"id", "player_id", "player_name", "room_code", "is_host", "role_id", "is_alive", "created_at", "last_seen"}).AddRow(sid, playerID, playerName, roomCode, isHost, sql.NullInt32{Int32: int32(roleID), Valid: true}, true, sql.NullTime{Time: createdAt, Valid: true}, sql.NullTime{Time: createdAt, Valid: true}))

	// Restore
	req2 := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	req2.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	c2.SetPath("/api/session")
	_ = sh.RestoreSession(c2)
	assert.Equal(t, http.StatusOK, rec2.Code)

	// Expect delete for delete
	fmt.Printf("DeleteSession args: sid: %T %v\n", sid, sid)
	mock.ExpectExec(`DELETE FROM sessions WHERE id = \$1`).
		WithArgs(sid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Delete
	req3 := httptest.NewRequest(http.MethodDelete, "/api/session", nil)
	req3.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rec3 := httptest.NewRecorder()
	c3 := e.NewContext(req3, rec3)
	c3.SetPath("/api/session")
	_ = sh.DeleteSession(c3)
	assert.Equal(t, http.StatusOK, rec3.Code)
}
