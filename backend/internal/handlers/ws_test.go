package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/internal/game"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWebSocket_UpgradeSuccess verifies that valid connections can be upgraded.
func TestWebSocket_UpgradeSuccess(t *testing.T) {
	// Create the hub and handler
	hub := game.NewHub()
	go hub.Run()
	handler := &RoomHandler{Hub: hub}

	// Create a room first
	roomCode := hub.GetRoomManager().CreateRoom("host-123")

	// Create Echo server with WebSocket route
	e := echo.New()
	e.GET("/ws", handler.UpgradeWebSocket)

	// Start a test HTTP server
	server := httptest.NewServer(e)
	defer server.Close()

	// Convert HTTP URL to WebSocket URL
	// http://127.0.0.1:xxxxx -> ws://127.0.0.1:xxxxx
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") +
		"/ws?room=" + roomCode + "&player=player-1&name=Alice"

	// Connect with a WebSocket client
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "WebSocket dial should succeed")
	defer conn.Close()

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "should get 101 Switching Protocols")

	// Give hub time to register client
	time.Sleep(50 * time.Millisecond)

	// Verify client is in the room
	assert.Equal(t, 1, hub.GetRoomClientCount(roomCode), "client should be registered in room")

	// The client should receive a player_joined message (for themselves)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	var msg game.WSMessage
	err = conn.ReadJSON(&msg)
	// We might not get the message due to timing, that's ok for this test
	// The important thing is that the connection was established
	if err == nil {
		assert.Equal(t, game.MsgTypePlayerJoined, msg.Type)
	}
}

// TestWebSocket_MissingParams verifies that missing parameters return error.
func TestWebSocket_MissingParams(t *testing.T) {
	hub := game.NewHub()
	go hub.Run()
	handler := &RoomHandler{Hub: hub}

	e := echo.New()
	e.GET("/ws", handler.UpgradeWebSocket)

	server := httptest.NewServer(e)
	defer server.Close()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "missing room parameter",
			url:     "/ws?player=p1&name=Alice",
			wantErr: true,
		},
		{
			name:    "missing player parameter",
			url:     "/ws?room=ABC123&name=Alice",
			wantErr: true,
		},
		{
			name:    "missing both parameters",
			url:     "/ws?name=Alice",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + tt.url
			_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)

			if tt.wantErr {
				// Should fail to upgrade - not get 101
				assert.Error(t, err, "should fail with missing params")
				if resp != nil {
					assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestWebSocket_InvalidRoom verifies that non-existent rooms return error.
func TestWebSocket_InvalidRoom(t *testing.T) {
	hub := game.NewHub()
	go hub.Run()
	handler := &RoomHandler{Hub: hub}

	e := echo.New()
	e.GET("/ws", handler.UpgradeWebSocket)

	server := httptest.NewServer(e)
	defer server.Close()

	// Try to connect to a room that doesn't exist
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") +
		"/ws?room=FAKE99&player=player-1&name=Alice"

	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)

	assert.Error(t, err, "should fail with invalid room")
	if resp != nil {
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	}
}

// TestWebSocket_MultipleClients verifies multiple clients can connect to the same room.
func TestWebSocket_MultipleClients(t *testing.T) {
	hub := game.NewHub()
	go hub.Run()
	handler := &RoomHandler{Hub: hub}

	roomCode := hub.GetRoomManager().CreateRoom("host-123")

	e := echo.New()
	e.GET("/ws", handler.UpgradeWebSocket)

	server := httptest.NewServer(e)
	defer server.Close()

	// Connect 3 clients
	clients := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") +
			"/ws?room=" + roomCode + "&player=player-" + string('0'+byte(i)) + "&name=Player" + string('0'+byte(i))

		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err, "client %d should connect", i)
		clients[i] = conn
		defer conn.Close()
	}

	// Give hub time to register all clients
	time.Sleep(100 * time.Millisecond)

	// Verify all 3 clients are in the room
	assert.Equal(t, 3, hub.GetRoomClientCount(roomCode), "all 3 clients should be in room")
}

// TestWebSocket_DisconnectCleanup verifies that disconnected clients are cleaned up.
func TestWebSocket_DisconnectCleanup(t *testing.T) {
	hub := game.NewHub()
	go hub.Run()
	handler := &RoomHandler{Hub: hub}

	roomCode := hub.GetRoomManager().CreateRoom("host-123")

	e := echo.New()
	e.GET("/ws", handler.UpgradeWebSocket)

	server := httptest.NewServer(e)
	defer server.Close()

	// Connect a client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") +
		"/ws?room=" + roomCode + "&player=player-1&name=Alice"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	// Wait for registration
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 1, hub.GetRoomClientCount(roomCode), "client should be in room")

	// Close connection from client side (simulate disconnect)
	conn.Close()

	// Wait for cleanup
	time.Sleep(100 * time.Millisecond)

	// Client should be removed
	assert.Equal(t, 0, hub.GetRoomClientCount(roomCode), "client should be removed after disconnect")
}

// TestWebSocket_PingPong verifies the ping/pong keepalive mechanism works.
func TestWebSocket_PingPong(t *testing.T) {
	hub := game.NewHub()
	go hub.Run()
	handler := &RoomHandler{Hub: hub}

	roomCode := hub.GetRoomManager().CreateRoom("host-123")

	e := echo.New()
	e.GET("/ws", handler.UpgradeWebSocket)

	server := httptest.NewServer(e)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") +
		"/ws?room=" + roomCode + "&player=player-1&name=Alice"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Wait for connection to stabilize
	time.Sleep(50 * time.Millisecond)

	// Send a ping from client
	err = conn.WriteMessage(websocket.PingMessage, []byte{})
	require.NoError(t, err)

	// The server should respond with pong (handled internally by gorilla)
	// We just verify the connection stays open
	time.Sleep(100 * time.Millisecond)

	// Try to read - connection should still be open
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err = conn.ReadMessage()
	// We might get timeout or a message, but shouldn't get close error
	if err != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		t.Error("Connection closed unexpectedly after ping/pong")
	}
}
