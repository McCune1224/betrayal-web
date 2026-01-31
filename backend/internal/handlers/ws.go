package handlers

import (
	"log"
	"net/http"

	"backend/internal/game"
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
	// Extract room code and player info from query parameters
	// URL format: /ws?room=CODE&player=ID&name=NAME
	roomCode := c.QueryParam("room")
	playerID := c.QueryParam("player")
	playerName := c.QueryParam("name")

	// Validate required parameters
	if roomCode == "" || playerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameters: room, player",
		})
	}

	// Validate room exists
	_, err := h.Hub.GetRoomManager().GetRoom(roomCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Room not found",
		})
	}

	// UPGRADE: Convert HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
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

	log.Printf("WebSocket client connected: room=%s player=%s", roomCode, playerID)

	return nil
}
