package handlers

import (
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"yourmodule/internal/game"
)

type RoomHandler struct {
	Hub *game.Hub
}

type CreateRoomRequest struct{}

type CreateRoomResponse struct {
	Code string `json:"code"`
	Host string `json:"host"`
}

type JoinRoomRequest struct {
	Name string `json:"name"`
}

type JoinRoomResponse struct {
	Code     string `json:"code"`
	PlayerID string `json:"player_id"`
	Phase    string `json:"phase"`
}

func (rh *RoomHandler) CreateRoom(c echo.Context) error {
	code := generateRoomCode()
	hostID := generatePlayerID()

	rh.Hub.GetRoomManager().CreateRoom(code, hostID)

	return c.JSON(http.StatusOK, CreateRoomResponse{
		Code: code,
		Host: hostID,
	})
}

func (rh *RoomHandler) JoinRoom(c echo.Context) error {
	code := c.Param("code")
	var req JoinRoomRequest
	if err := c.BindJSON(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	room, player, err := rh.Hub.GetRoomManager().JoinRoom(code, req.Name)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, JoinRoomResponse{
		Code:     code,
		PlayerID: player.ID,
		Phase:    room.Phase,
	})
}

func (rh *RoomHandler) UpgradeWebSocket(c echo.Context) error {
	// TODO: Implement WebSocket upgrade
	return c.String(http.StatusOK, "WebSocket not yet implemented")
}

func generateRoomCode() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generatePlayerID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
