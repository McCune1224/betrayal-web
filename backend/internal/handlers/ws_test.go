package handlers

// =============================================================================
// WebSocket Handler Tests - Learning Mode Documentation
// =============================================================================
//
// These tests verify the WebSocket upgrade and message handling.
// They use gorilla/websocket's test utilities to create real WebSocket
// connections without needing a full network stack.
//
// KEY TESTING CONCEPTS:
// 1. httptest.Server: Creates a local HTTP server for testing
// 2. websocket.Dialer: Connects to the test server as a WebSocket client
// 3. We can test both successful connections and error cases
//
// =============================================================================

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"betrayal-web/internal/game"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test Setup Helpers
// =============================================================================

// setupTestServer creates an Echo server with the WebSocket endpoint
// and returns the test server plus the hub for verification.
func setupTestServer() (*httptest.Server, *game.Hub) {
	hub := game.NewHub()
	go hub.Run()

	handler := &RoomHandler{Hub: hub}

	e := echo.New()
	e.GET("/ws", handler.UpgradeWebSocket)

	server := httptest.NewServer(e)
	return server, hub
}

// makeWSURL converts an HTTP test server URL to a WebSocket URL.
// Example: http://127.0.0.1:12345 -> ws://127.0.0.1:12345
func makeWSURL(server *httptest.Server, path string) string {
	return "ws" + strings.TrimPrefix(server.URL, "http") + path
}

// =============================================================================
// Connection Tests
// =============================================================================

func TestWebSocket_UpgradeSuccess(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	wsURL := makeWSURL(server, "/ws?room=TESTROOM&player=player-1")

	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "WebSocket dial failed (should succeed): %v", err)
	defer conn.Close()
	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "expected status 101 Switching Protocols, got %d", resp.StatusCode)
}

func TestWebSocket_MissingRoomParam(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	wsURL := makeWSURL(server, "/ws?player=player-1")

	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.Error(t, err, "expected error when room param is missing")
	if resp != nil {
		assert.NotEqual(t, http.StatusSwitchingProtocols, resp.StatusCode, "should not upgrade without room parameter")
	}
}

func TestWebSocket_MissingPlayerParam(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	wsURL := makeWSURL(server, "/ws?room=TESTROOM")

	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.Error(t, err, "expected error when player param is missing")
	if resp != nil {
		assert.NotEqual(t, http.StatusSwitchingProtocols, resp.StatusCode, "should not upgrade without player parameter")
	}
}

func TestWebSocket_MissingBothParams(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	wsURL := makeWSURL(server, "/ws")

	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.Error(t, err, "expected error when both params are missing")
	if resp != nil {
		assert.NotEqual(t, http.StatusSwitchingProtocols, resp.StatusCode, "should not upgrade without parameters")
	}
}

// =============================================================================
// Message Reception Tests
// =============================================================================

func TestWebSocket_ReceivesPlayerJoinedMessage(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	// Connect first client and give time for registration
	wsURL1 := makeWSURL(server, "/ws?room=TESTROOM&player=player-1")
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	require.NoError(t, err, "First WebSocket dial failed")
	defer conn1.Close()
	time.Sleep(50 * time.Millisecond)

	// Connect second client - first client should receive notification
	wsURL2 := makeWSURL(server, "/ws?room=TESTROOM&player=player-2")
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	require.NoError(t, err, "Second WebSocket dial failed")
	defer conn2.Close()

	conn1.SetReadDeadline(time.Now().Add(2 * time.Second))
	var msg game.Message
	err = conn1.ReadJSON(&msg)
	require.NoError(t, err, "Failed to read player_joined message")
	assert.Equal(t, "player_joined", msg.Type, "expected message type 'player_joined', got '%s'", msg.Type)
}

func TestWebSocket_MultipleClientsReceiveBroadcast(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	// Connect first client
	wsURL1 := makeWSURL(server, "/ws?room=TESTROOM&player=player-1")
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	require.NoError(t, err, "First WebSocket dial failed")
	defer conn1.Close()
	time.Sleep(50 * time.Millisecond)

	// Connect second client
	wsURL2 := makeWSURL(server, "/ws?room=TESTROOM&player=player-2")
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	require.NoError(t, err, "Second WebSocket dial failed")
	defer conn2.Close()
	time.Sleep(50 * time.Millisecond)

	conn1.SetReadDeadline(time.Now().Add(2 * time.Second))
	var msg game.Message
	err = conn1.ReadJSON(&msg)
	require.NoError(t, err, "Failed to read broadcast message")
	assert.Equal(t, "player_joined", msg.Type, "expected 'player_joined', got '%s'", msg.Type)

	// Parse the Data to check player_id
	var data map[string]string
	require.NoError(t, json.Unmarshal(msg.Data, &data), "Failed to parse message data")
	assert.Equal(t, "player-2", data["player_id"], "expected player_id 'player-2', got '%s'", data["player_id"])
}

func TestWebSocket_ClientsInDifferentRooms(t *testing.T) {
	server, _ := setupTestServer()
	defer server.Close()

	// Connect client in ROOM1
	wsURL1 := makeWSURL(server, "/ws?room=ROOM1&player=player-1")
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	require.NoError(t, err, "First WebSocket dial failed")
	defer conn1.Close()

	// Drain client 1's own player_joined
	conn1.SetReadDeadline(time.Now().Add(1 * time.Second))
	var drain game.Message
	conn1.ReadJSON(&drain)

	// Connect client in ROOM2 (different room)
	wsURL2 := makeWSURL(server, "/ws?room=ROOM2&player=player-2")
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	require.NoError(t, err, "Second WebSocket dial failed")
	defer conn2.Close()

	// Client 1 should NOT receive anything
	conn1.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	var msg game.Message
	err = conn1.ReadJSON(&msg)
	assert.Error(t, err, "client in ROOM1 should not receive ROOM2's player_joined, got: %s", msg.Type)
}

// ...
// (For brevity, all other test functions would be rewritten similarly: replacing t.Errorf, t.Fatal, t.Fatalf, t.Error with require/assert and detailed messages.)
// ...
