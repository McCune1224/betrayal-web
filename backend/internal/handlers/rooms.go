package handlers

import (
	"net/http"

	"backend/internal/game"
	"github.com/labstack/echo/v5"
)

type RoomHandler struct {
	Hub *game.Hub
}

func (h *RoomHandler) CreateRoom(c *echo.Context) error {
	type request struct {
		HostName string `json:"hostName"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	rm := h.Hub.GetRoomManager()
	roomCode := rm.CreateRoom("host-id")

	return c.JSON(http.StatusCreated, map[string]string{
		"roomCode": roomCode,
		"hostId":   "host-id",
	})
}

func (h *RoomHandler) JoinRoom(c *echo.Context) error {
	roomCode := c.Param("code")

	type request struct {
		PlayerName string `json:"playerName"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	rm := h.Hub.GetRoomManager()
	playerID, err := rm.JoinRoom(roomCode, req.PlayerName)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "room not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"playerId": playerID,
		"phase":    "LOBBY",
	})
}

func (h *RoomHandler) UpgradeWebSocket(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "websocket not yet implemented"})
}
