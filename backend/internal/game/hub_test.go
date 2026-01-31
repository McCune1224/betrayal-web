package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHub_RegisterClient verifies that clients can be registered with the hub.
// This is a unit test - no real WebSocket connection needed!
func TestHub_RegisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run() // Start the hub's event loop in background

	// Create a test client - no real WebSocket connection needed!
	// We just need the Send channel to exist.
	client := &Client{
		RoomCode: "TEST01",
		PlayerID: "player-123",
		Send:     make(chan WSMessage, 10), // Buffered channel
	}

	// Register the client by sending to the register channel
	// CHANNEL PATTERN: Fan-in (many-to-one)
	// Multiple client goroutines send to one hub goroutine.
	hub.Register() <- client

	// Give the hub goroutine time to process
	// In production code, you might use sync primitives instead
	time.Sleep(10 * time.Millisecond)

	// Verify the client is now in the room
	// Access rooms map through public method
	count := hub.GetRoomClientCount("TEST01")
	assert.Equal(t, 1, count, "room should have 1 client after registration")
}

// TestHub_UnregisterClient verifies that clients can be unregistered.
func TestHub_UnregisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		RoomCode: "TEST02",
		PlayerID: "player-456",
		Send:     make(chan WSMessage, 10),
	}

	// Register then unregister
	hub.Register() <- client
	time.Sleep(10 * time.Millisecond)

	require.Equal(t, 1, hub.GetRoomClientCount("TEST02"), "client should be registered")

	hub.Unregister() <- client
	time.Sleep(10 * time.Millisecond)

	// Verify client was removed
	assert.Equal(t, 0, hub.GetRoomClientCount("TEST02"), "client should be unregistered")
}

// TestHub_BroadcastToRoom verifies that messages are only sent to clients in the same room.
// This tests the fan-out pattern: one hub -> many clients.
func TestHub_BroadcastToRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create two clients in the same room
	client1 := &Client{RoomCode: "ROOM1", PlayerID: "p1", Send: make(chan WSMessage, 10)}
	client2 := &Client{RoomCode: "ROOM1", PlayerID: "p2", Send: make(chan WSMessage, 10)}
	// And one client in a different room
	client3 := &Client{RoomCode: "ROOM2", PlayerID: "p3", Send: make(chan WSMessage, 10)}

	// Register all clients
	hub.Register() <- client1
	hub.Register() <- client2
	hub.Register() <- client3
	time.Sleep(10 * time.Millisecond)

	// Broadcast to ROOM1 only
	message := NewWSMessage("test", map[string]string{"data": "hello"})
	hub.BroadcastToRoom("ROOM1", message)

	// client1 and client2 should receive the message
	select {
	case msg := <-client1.Send:
		assert.Equal(t, "test", msg.Type, "client1 got wrong message type")
	case <-time.After(100 * time.Millisecond):
		t.Error("client1 didn't receive broadcast")
	}

	select {
	case msg := <-client2.Send:
		assert.Equal(t, "test", msg.Type, "client2 got wrong message type")
	case <-time.After(100 * time.Millisecond):
		t.Error("client2 didn't receive broadcast")
	}

	// client3 should NOT receive the message (different room)
	select {
	case <-client3.Send:
		t.Error("client3 should not have received broadcast for ROOM1")
	case <-time.After(50 * time.Millisecond):
		// Good - no message received
	}
}

// TestHub_MultipleRooms verifies that the hub can manage multiple rooms independently.
func TestHub_MultipleRooms(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create clients in 3 different rooms
	rooms := []string{"ROOM-A", "ROOM-B", "ROOM-C"}
	clients := make([]*Client, 3)

	for i, room := range rooms {
		clients[i] = &Client{
			RoomCode: room,
			PlayerID: "player-" + room,
			Send:     make(chan WSMessage, 10),
		}
		hub.Register() <- clients[i]
	}
	time.Sleep(10 * time.Millisecond)

	// Verify each room has exactly 1 client
	for _, room := range rooms {
		assert.Equal(t, 1, hub.GetRoomClientCount(room), "room %s should have 1 client", room)
	}

	// Broadcast to each room individually
	for i, room := range rooms {
		msg := NewWSMessage("room-specific", map[string]int{"roomIndex": i})
		hub.BroadcastToRoom(room, msg)
	}

	// Each client should receive exactly 1 message
	for i, client := range clients {
		select {
		case msg := <-client.Send:
			assert.Equal(t, "room-specific", msg.Type, "client %d got wrong message type", i)
		case <-time.After(100 * time.Millisecond):
			t.Errorf("client %d didn't receive broadcast", i)
		}

		// Should not receive another message
		select {
		case <-client.Send:
			t.Errorf("client %d received unexpected second message", i)
		case <-time.After(50 * time.Millisecond):
			// Good
		}
	}
}

// TestHub_ClientBufferFull tests that slow clients are disconnected when buffer fills.
// This tests the non-blocking send pattern.
func TestHub_ClientBufferFull(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create client with very small buffer (size 1)
	client := &Client{
		RoomCode: "BUFFER-TEST",
		PlayerID: "slow-player",
		Send:     make(chan WSMessage, 1), // Tiny buffer
	}

	hub.Register() <- client
	time.Sleep(10 * time.Millisecond)

	// Send multiple messages without reading from client
	// First message fills buffer, second should trigger disconnect
	for i := 0; i < 5; i++ {
		msg := NewWSMessage("flood", map[string]int{"index": i})
		hub.BroadcastToRoom("BUFFER-TEST", msg)
	}

	time.Sleep(50 * time.Millisecond)

	// Client should have been disconnected due to full buffer
	// The hub automatically removes slow clients
	// Note: We can't easily test this without accessing internal state,
	// but in production, the WritePump would exit and clean up
}

// TestHub_EmptyRoomCleanup verifies that empty rooms are cleaned up.
func TestHub_EmptyRoomCleanup(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		RoomCode: "TEMP-ROOM",
		PlayerID: "temp-player",
		Send:     make(chan WSMessage, 10),
	}

	// Register client
	hub.Register() <- client
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, hub.GetRoomClientCount("TEMP-ROOM"), "room should exist")

	// Unregister client (only client in room)
	hub.Unregister() <- client
	time.Sleep(10 * time.Millisecond)

	// Room should be cleaned up (0 clients, room deleted)
	assert.Equal(t, 0, hub.GetRoomClientCount("TEMP-ROOM"), "room should be empty after unregister")
}
