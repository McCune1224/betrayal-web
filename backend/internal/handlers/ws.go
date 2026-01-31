package handlers

import (
	"net/http"

	"backend/internal/game"
	"backend/internal/logging"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

// upgrader configures the WebSocket upgrade from HTTP to WebSocket protocol.
// CheckOrigin allows connections from any origin - suitable for development.
// In production, you should restrict this to your domain.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// UpgradeWebSocket upgrades an HTTP connection to WebSocket and starts the client.
// This is the entry point for WebSocket connections.
func (h *RoomHandler) UpgradeWebSocket(c *echo.Context) error {
	ctx := (*c).Request().Context()
	logger := logging.WithContext(ctx)

	// Extract room code and player info from query parameters
	// URL format: /ws?room=CODE&player=ID&name=NAME
	roomCode := (*c).QueryParam("room")
	playerID := (*c).QueryParam("player")
	playerName := (*c).QueryParam("name")

	logger.Debug("websocket_upgrade_attempt",
		"room_code", roomCode,
		"player_id", playerID,
		"player_name", playerName,
	)

	// Validate required parameters
	if roomCode == "" || playerID == "" {
		logger.Warn("websocket_upgrade_failed",
			"reason", "missing_parameters",
			"room_code", roomCode,
			"player_id", playerID,
		)
		return (*c).JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameters: room, player",
		})
	}

	// Validate room exists
	_, err := h.Hub.GetRoomManager().GetRoom(roomCode)
	if err != nil {
		logger.Warn("websocket_upgrade_failed",
			"reason", "room_not_found",
			"room_code", roomCode,
		)
		return (*c).JSON(http.StatusNotFound, map[string]string{
			"error": "Room not found",
		})
	}

	// UPGRADE: Convert HTTP connection to WebSocket
	conn, err := upgrader.Upgrade((*c).Response(), (*c).Request(), nil)
	if err != nil {
		logger.Error("websocket_upgrade_failed",
			"reason", "upgrade_error",
			"room_code", roomCode,
			"player_id", playerID,
			"error", err,
		)
		return err
	}

	// Create new client instance
	client := &game.Client{
		Hub:        h.Hub,
		Conn:       conn,
		Send:       make(chan game.WSMessage, 256),
		RoomCode:   roomCode,
		PlayerID:   playerID,
		PlayerName: playerName,
		IsHost:     false,
	}

	// Register client with the hub
	h.Hub.Register() <- client

	// Start read and write goroutines
	go client.WritePump()
	go client.ReadPump()

	// Broadcast to room that player joined
	joinMsg := game.NewWSMessage(game.MsgTypePlayerJoined, game.PlayerJoinedData{
		PlayerID:   playerID,
		PlayerName: playerName,
		IsHost:     client.IsHost,
	})
	h.Hub.BroadcastToRoom(roomCode, joinMsg)

	logger.Info("websocket_client_connected",
		"room_code", roomCode,
		"player_id", playerID,
		"player_name", playerName,
		"remote_addr", (*c).Request().RemoteAddr,
	)

	return nil
}
