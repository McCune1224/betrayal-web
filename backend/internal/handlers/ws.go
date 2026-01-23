package handlers

// =============================================================================
// WebSocket Handler - Learning Mode Documentation
// =============================================================================
//
// This file implements the WebSocket upgrade and message handling for real-time
// game communication. If you're new to WebSockets in Go, read the comments
// carefully - they explain the "why" behind each pattern.
//
// KEY CONCEPTS:
// 1. HTTP Upgrade: WebSocket starts as HTTP, then "upgrades" to a persistent connection
// 2. Goroutines: Each client gets TWO goroutines (read + write) for concurrent I/O
// 3. Channels: We use channels to communicate between goroutines safely
// 4. Hub Pattern: A central hub manages all connections and routes messages
//
// =============================================================================

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"betrayal-web/internal/game"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// =============================================================================
// Constants - Tuning Parameters for WebSocket Behavior
// =============================================================================

const (
	// writeWait is the maximum time allowed to write a message to the client.
	// If the write takes longer than this, we assume the connection is dead.
	// This prevents goroutines from hanging forever on slow/dead connections.
	writeWait = 10 * time.Second

	// pongWait is how long we wait for a "pong" response from the client.
	// WebSocket has a built-in ping/pong mechanism for detecting dead connections.
	// If the client doesn't respond within this time, we close the connection.
	pongWait = 60 * time.Second

	// pingPeriod is how often we send "ping" messages to the client.
	// This must be less than pongWait, so we send pings before the deadline.
	// Formula: pingPeriod < pongWait (we use 90% of pongWait)
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize is the maximum size of incoming messages in bytes.
	// This prevents clients from sending huge messages that could crash the server.
	// 512KB should be more than enough for game messages.
	maxMessageSize = 512 * 1024
)

// =============================================================================
// WebSocket Upgrader Configuration
// =============================================================================

// upgrader is the gorilla/websocket upgrader that converts HTTP to WebSocket.
// It's configured once and reused for all connections.
var upgrader = websocket.Upgrader{
	// ReadBufferSize and WriteBufferSize control the I/O buffer sizes.
	// Larger buffers = more memory per connection but better throughput.
	// 1KB is fine for our small JSON messages.
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// CheckOrigin controls which origins (domains) can connect.
	// In production, you should validate this properly!
	// For development, we allow all origins.
	//
	// SECURITY NOTE: In production, replace this with:
	//   CheckOrigin: func(r *http.Request) bool {
	//       origin := r.Header.Get("Origin")
	//       return origin == "https://yourdomain.com"
	//   }
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// =============================================================================
// UpgradeWebSocket - The Entry Point
// =============================================================================

// UpgradeWebSocket handles the GET /ws endpoint.
// It upgrades an HTTP connection to a WebSocket and starts the read/write pumps.
//
// Query Parameters:
//   - room: The room code to join (required)
//   - player: The player ID (required)
//
// Example: GET /ws?room=ABC123&player=def456
func (rh *RoomHandler) UpgradeWebSocket(c echo.Context) error {
	// STEP 1: Extract and validate query parameters
	// The frontend must provide the room code and player ID in the URL.
	roomCode := c.QueryParam("room")
	playerID := c.QueryParam("player")

	if roomCode == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "room parameter is required",
		})
	}
	if playerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "player parameter is required",
		})
	}

	// STEP 2: Upgrade HTTP to WebSocket
	// This is where the magic happens. The upgrader:
	// 1. Validates the HTTP request has the right headers
	// 2. Sends back a "101 Switching Protocols" response
	// 3. Returns a *websocket.Conn for bidirectional communication
	//
	// After this call, we're no longer speaking HTTP - it's pure WebSocket.
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		// If upgrade fails, the client probably sent a bad request.
		// The upgrader already wrote an error response, so just log and return.
		log.Printf("WebSocket upgrade failed: %v", err)
		return nil // Don't write another response
	}

	// STEP 3: Create a Client struct to represent this connection
	// The Client ties together:
	// - The WebSocket connection (for reading/writing)
	// - The room/player info (for routing messages)
	// - A Send channel (for queuing outgoing messages)
	client := &game.Client{
		RoomCode: roomCode,
		PlayerID: playerID,
		Conn:     conn,
		Hub:      rh.Hub,
		// BUFFERED CHANNEL: Size 256 means we can queue up to 256 messages
		// before the sender blocks. This prevents slow clients from blocking
		// the hub's broadcast loop.
		Send: make(chan game.Message, 256),
	}

	// STEP 4: Register the client with the Hub
	// This is a BLOCKING send to the hub's register channel.
	// The hub's Run() goroutine will receive this and add the client
	// to the appropriate room's connection list.
	rh.Hub.Register(client)

	// STEP 5: Start the read and write pumps in separate goroutines
	//
	// WHY TWO GOROUTINES?
	// WebSocket connections are full-duplex (can read and write simultaneously),
	// but the operations themselves are blocking:
	// - ReadMessage() blocks until a message arrives
	// - WriteMessage() blocks until the message is sent
	//
	// If we used one goroutine, we couldn't read while writing (or vice versa).
	// Two goroutines allow truly concurrent I/O.
	//
	// WHY GOROUTINES AT ALL?
	// This HTTP handler must return quickly so Echo can handle other requests.
	// The goroutines keep running after this function returns, maintaining
	// the WebSocket connection for as long as needed.
	go writePump(client)
	go readPump(client, rh.Hub)

	// STEP 6: Notify other players that someone joined
	// Broadcast a "player_joined" message to everyone in the room.
	rh.Hub.BroadcastToRoom(roomCode, game.Message{
		Type: "player_joined",
		Data: mustMarshal(map[string]string{
			"player_id": playerID,
		}),
	})

	// Return nil because the upgrader already sent the response.
	// The connection is now handled by the goroutines.
	return nil
}

// =============================================================================
// readPump - Reads Messages from the Browser
// =============================================================================

// readPump reads messages from the WebSocket connection and processes them.
// It runs in its own goroutine, one per client.
//
// LIFECYCLE:
// 1. Set up read deadline and pong handler
// 2. Loop forever, reading messages
// 3. When read fails (error or close), clean up and exit
//
// WHY "PUMP"?
// It continuously "pumps" messages from the network into the application,
// like a water pump moves water. It's a common Go naming convention.
func readPump(client *game.Client, hub *game.Hub) {
	// DEFER: This cleanup runs when the function exits (for any reason).
	// It's CRITICAL for resource management - ensures we always clean up.
	//
	// Order matters! Defers run in LIFO (last-in-first-out) order:
	// 1. First, close the connection (frees OS resources)
	// 2. Then, unregister from hub (removes from room, closes Send channel)
	//
	// We put unregister first in the code so it runs LAST.
	defer func() {
		// Tell the hub we're leaving. This:
		// - Removes us from the room's client list
		// - Closes our Send channel (signals writePump to exit)
		// - Cleans up empty rooms
		hub.Unregister(client)

		// Close the WebSocket connection.
		// This is safe to call multiple times.
		client.Conn.Close()
	}()

	// Configure the connection for reading
	//
	// SetReadLimit: Reject messages larger than maxMessageSize.
	// This prevents malicious clients from sending huge messages.
	client.Conn.SetReadLimit(maxMessageSize)

	// SetReadDeadline: If we don't receive ANYTHING within pongWait,
	// the next read will fail with a timeout error.
	// This is reset every time we receive a pong.
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// SetPongHandler: Called when we receive a "pong" from the browser.
	// We use this to reset the read deadline - the connection is still alive.
	//
	// HOW PING/PONG WORKS:
	// 1. writePump sends a "ping" every pingPeriod (54 seconds)
	// 2. Browser automatically responds with "pong"
	// 3. This handler receives the pong and resets the deadline
	// 4. If no pong arrives within pongWait (60 seconds), connection is dead
	client.Conn.SetPongHandler(func(string) error {
		// Connection is alive! Reset the deadline.
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// THE READ LOOP
	// This loops forever until an error occurs (connection closed, timeout, etc.)
	for {
		// ReadMessage blocks until a message arrives.
		// Returns:
		// - messageType: Text (1) or Binary (2) - we expect Text (JSON)
		// - message: The raw bytes of the message
		// - err: Non-nil if something went wrong
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			// ANY error means the connection is dead or broken.
			// Common causes:
			// - Client closed the browser tab
			// - Network disconnected
			// - Read deadline exceeded (no pong received)
			// - Client sent invalid WebSocket frame
			//
			// We don't distinguish between error types - just exit.
			// The defer will clean up.
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse the JSON message
		var msg game.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			// Invalid JSON - log and continue (don't crash the connection)
			log.Printf("Invalid JSON from client %s: %v", client.PlayerID, err)
			continue
		}

		// Handle the message based on its type
		handleMessage(client, hub, msg)
	}
}

// =============================================================================
// writePump - Writes Messages to the Browser
// =============================================================================

// writePump writes messages from the Send channel to the WebSocket connection.
// It runs in its own goroutine, one per client.
//
// WHY A SEPARATE GOROUTINE?
// WebSocket connections are NOT thread-safe for concurrent writes.
// If two goroutines try to write at the same time, the connection corrupts.
// This goroutine "owns" all writes to this connection.
//
// COMMUNICATION:
// Other parts of the code (hub broadcasts, direct sends) put messages
// on the client.Send channel. This goroutine reads from that channel
// and writes to the actual network connection.
func writePump(client *game.Client) {
	// Ticker sends a signal every pingPeriod (54 seconds).
	// We use this to send ping messages to keep the connection alive.
	ticker := time.NewTicker(pingPeriod)

	// DEFER: Clean up when the function exits
	defer func() {
		// Stop the ticker to free the timer resources
		ticker.Stop()
		// Close the connection
		client.Conn.Close()
	}()

	// THE WRITE LOOP
	for {
		// SELECT: Wait for one of these things to happen:
		// 1. A message arrives on client.Send
		// 2. The ticker fires (time to send a ping)
		select {
		case message, ok := <-client.Send:
			// Set a deadline for this write operation.
			// If the write takes longer than writeWait, it fails.
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// THE CHANNEL WAS CLOSED!
				// This happens when:
				// - The hub unregisters us (another client saw us disconnect)
				// - The readPump called hub.Unregister
				//
				// Send a WebSocket close message to the browser and exit.
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write the message as JSON to the browser
			err := client.Conn.WriteJSON(message)
			if err != nil {
				// Write failed - connection is dead
				// The defer will close the connection
				return
			}

		case <-ticker.C:
			// TIME TO SEND A PING!
			// This keeps the connection alive and detects dead connections.
			//
			// Set a write deadline for the ping
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			// Send a WebSocket ping frame
			// The browser will automatically respond with a pong
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				// Couldn't send ping - connection is dead
				return
			}
		}
	}
}

// =============================================================================
// Message Handling
// =============================================================================

// handleMessage processes an incoming message based on its type.
// This is where you add game-specific logic.
func handleMessage(client *game.Client, hub *game.Hub, msg game.Message) {
	// Route the message based on its type
	switch msg.Type {
	case "chat":
		// Broadcast chat messages to everyone in the room
		hub.BroadcastToRoom(client.RoomCode, game.Message{
			Type: "chat",
			Data: mustMarshal(map[string]interface{}{
				"player_id": client.PlayerID,
				"message":   string(msg.Data),
			}),
		})

	case "action":
		// Player submitted a game action
		// TODO: Validate the action and update game state
		log.Printf("Action from %s: %s", client.PlayerID, string(msg.Data))

	case "advance_phase":
		// Host wants to advance the game phase
		// TODO: Verify this is the host, then advance
		log.Printf("Advance phase request from %s", client.PlayerID)

	default:
		log.Printf("Unknown message type from %s: %s", client.PlayerID, msg.Type)
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// mustMarshal converts a value to JSON, panicking on error.
// Only use this for values you KNOW will marshal successfully.
func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
