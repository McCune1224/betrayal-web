package game

import (
	"time"

	"backend/internal/logging"
	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
// Each connected player gets one Client instance.
// This type is defined in game package so Hub can reference it.
type Client struct {
	// Hub is the central coordinator that manages all clients and rooms.
	// WHY: Hub provides thread-safe access to shared state (rooms map).
	Hub *Hub

	// Conn is the underlying WebSocket connection from gorilla/websocket.
	// WHY: This is our communication channel with the browser.
	// Note: This is capitalized so handlers package can access it.
	Conn *websocket.Conn

	// Send is a buffered channel for outgoing messages.
	// WHY BUFFERED: Prevents slow clients from blocking the hub.
	// If a client is slow to read, messages accumulate here instead of blocking broadcasts.
	Send chan WSMessage

	// RoomCode identifies which game room this client belongs to.
	// WHY: Used to route messages only to clients in the same room.
	RoomCode string

	// PlayerID uniquely identifies the player across reconnects.
	// WHY: Allows rejoining after disconnect without losing state.
	PlayerID string

	// PlayerName is the display name shown to other players.
	PlayerName string

	// IsHost indicates if this client is the room host.
	// WHY: Host has special permissions (start game, kick players, etc.).
	IsHost bool

	// connectedAt tracks when the client connected for metrics
	connectedAt time.Time
}

// Timeouts for WebSocket keepalive mechanism
// These values prevent half-open connections and detect dead clients
const (
	// writeWait is the time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// pongWait is the time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// pingPeriod is how often to send ping frames.
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize is the maximum allowed message size from client.
	maxMessageSize = 512
)

// ReadPump pumps messages from the WebSocket connection to the hub.
// WHY GOROUTINE: Each client needs independent read loop.
// The read loop blocks waiting for messages from the browser.
// If we didn't use a goroutine, one slow client would block all others.
func (c *Client) ReadPump() {
	logger := logging.WSLogger(c.RoomCode, c.PlayerID)
	c.connectedAt = time.Now()

	// CLEANUP: Always use defer for WebSocket cleanup.
	// This runs when the function returns, even if it panics.
	defer func() {
		// 1. Tell the hub we're leaving (removes us from room)
		c.Hub.Unregister() <- c
		// 2. Close the network connection (frees OS resources)
		c.Conn.Close()

		duration := time.Since(c.connectedAt)
		logger.Info("client_disconnected",
			"connection_duration_ms", duration.Milliseconds(),
		)
	}()

	// Set maximum message size to prevent memory exhaustion attacks
	c.Conn.SetReadLimit(maxMessageSize)

	// Set read deadline - detects dead connections
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// Set pong handler - resets deadline on pong response
	c.Conn.SetPongHandler(func(string) error {
		logger.Debug("pong_received")
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Read loop: continuously read messages from the browser
	for {
		var msg WSMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Warn("websocket_error",
					"error", err,
					"error_type", "unexpected_close",
				)
			} else {
				logger.Debug("websocket_closed",
					"error", err,
				)
			}
			break
		}

		c.handleIncomingMessage(msg)
	}
}

// handleIncomingMessage routes incoming messages to the appropriate handler
func (c *Client) handleIncomingMessage(msg WSMessage) {
	logger := logging.WSLogger(c.RoomCode, c.PlayerID)

	logger.Debug("message_received",
		"message_type", msg.Type,
	)

	switch msg.Type {
	case "join_room":
		c.handleJoinRoom(msg)

	case "submit_action":
		c.handleSubmitAction(msg)

	case "host_command":
		if !c.IsHost {
			logger.Warn("unauthorized_host_command",
				"player_id", c.PlayerID,
			)
			c.sendError("not_host", "Only the host can perform this action")
			return
		}
		c.handleHostCommand(msg)

	case "ping":
		c.Send <- NewWSMessage("pong", nil)

	default:
		logger.Warn("unknown_message_type",
			"message_type", msg.Type,
		)
		c.sendError("unknown_message_type", "Unknown message type: "+msg.Type)
	}
}

// handleJoinRoom processes a join room message (typically for rejoins)
func (c *Client) handleJoinRoom(msg WSMessage) {
	logger := logging.WSLogger(c.RoomCode, c.PlayerID)

	data := PlayerJoinedData{
		PlayerID:   c.PlayerID,
		PlayerName: c.PlayerName,
		IsHost:     c.IsHost,
	}
	broadcast := NewWSMessage(MsgTypePlayerRejoined, data)
	c.Hub.BroadcastToRoom(c.RoomCode, broadcast)

	logger.Info("player_rejoined",
		"player_name", c.PlayerName,
		"is_host", c.IsHost,
	)
}

// handleSubmitAction processes action submissions from players
func (c *Client) handleSubmitAction(msg WSMessage) {
	logger := logging.WSLogger(c.RoomCode, c.PlayerID)

	logger.Info("action_submitted",
		"player_name", c.PlayerName,
		"action_data", msg.Data,
	)
}

// handleHostCommand processes host-only commands
func (c *Client) handleHostCommand(msg WSMessage) {
	logger := logging.WSLogger(c.RoomCode, c.PlayerID)

	logger.Info("host_command",
		"player_name", c.PlayerName,
		"command_data", msg.Data,
	)
}

// sendError sends an error message to this client only
func (c *Client) sendError(code, message string) {
	logger := logging.WSLogger(c.RoomCode, c.PlayerID)

	logger.Debug("sending_error_to_client",
		"error_code", code,
		"error_message", message,
	)

	errData := ErrorData{
		Code:    code,
		Message: message,
	}
	c.Send <- NewWSMessage(MsgTypeError, errData)
}

// WritePump pumps messages from the hub to the WebSocket connection.
// WHY SEPARATE WRITE GOROUTINE: WebSocket connections are NOT thread-safe
// for concurrent writes. This goroutine owns all writes to this connection.
func (c *Client) WritePump() {
	logger := logging.WSLogger(c.RoomCode, c.PlayerID)
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
		logger.Debug("write_pump_stopped")
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// Channel was closed by hub
				logger.Debug("send_channel_closed")
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteJSON(message)
			if err != nil {
				logger.Warn("write_error",
					"error", err,
				)
				return
			}

			logger.Debug("message_sent",
				"message_type", message.Type,
			)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Warn("ping_error",
					"error", err,
				)
				return
			}
			logger.Debug("ping_sent")
		}
	}
}
