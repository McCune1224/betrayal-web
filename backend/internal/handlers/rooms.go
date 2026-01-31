package handlers

import (
	"net/http"

	"backend/internal/game"
	"backend/internal/logging"
	"github.com/labstack/echo/v5"
)

type RoomHandler struct {
	Hub *game.Hub
}

func (h *RoomHandler) CreateRoom(c *echo.Context) error {
	ctx := (*c).Request().Context()
	logger := logging.WithContext(ctx)

	type request struct {
		HostName string `json:"hostName"`
	}

	var req request
	if err := (*c).Bind(&req); err != nil {
		logger.Warn("create_room_failed",
			"reason", "invalid_request",
			"error", err,
		)
		return (*c).JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	logger.Debug("creating_room",
		"host_name", req.HostName,
	)

	rm := h.Hub.GetRoomManager()
	roomCode := rm.CreateRoom("host-id")

	logger.Info("room_created",
		"room_code", roomCode,
		"host_name", req.HostName,
	)

	return (*c).JSON(http.StatusCreated, map[string]string{
		"roomCode": roomCode,
		"hostId":   "host-id",
	})
}

func (h *RoomHandler) JoinRoom(c *echo.Context) error {
	ctx := (*c).Request().Context()
	logger := logging.WithContext(ctx)

	roomCode := (*c).Param("code")

	type request struct {
		PlayerName string `json:"playerName"`
	}

	var req request
	if err := (*c).Bind(&req); err != nil {
		logger.Warn("join_room_failed",
			"reason", "invalid_request",
			"room_code", roomCode,
			"error", err,
		)
		return (*c).JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	logger.Debug("joining_room",
		"room_code", roomCode,
		"player_name", req.PlayerName,
	)

	rm := h.Hub.GetRoomManager()
	playerID, err := rm.JoinRoom(roomCode, req.PlayerName)
	if err != nil {
		logger.Warn("join_room_failed",
			"reason", "room_not_found",
			"room_code", roomCode,
			"player_name", req.PlayerName,
		)
		return (*c).JSON(http.StatusNotFound, map[string]string{"error": "room not found"})
	}

	logger.Info("player_joined",
		"room_code", roomCode,
		"player_id", playerID,
		"player_name", req.PlayerName,
	)

	return (*c).JSON(http.StatusOK, map[string]string{
		"playerId": playerID,
		"phase":    "LOBBY",
	})
}
