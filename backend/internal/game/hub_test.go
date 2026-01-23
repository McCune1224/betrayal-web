package game

import (
	"testing"
	"time"
)

// =============================================================================
// Hub Constructor Tests
// =============================================================================

func TestNewHub(t *testing.T) {
	// NewHub should create a Hub with initialized channels and RoomManager
	hub := NewHub()

	if hub == nil {
		t.Fatal("NewHub returned nil")
	}
	if hub.rooms == nil {
		t.Error("rooms map should be initialized")
	}
	if hub.register == nil {
		t.Error("register channel should be initialized")
	}
	if hub.unregister == nil {
		t.Error("unregister channel should be initialized")
	}
	if hub.roomManager == nil {
		t.Error("roomManager should be initialized")
	}
}

func TestHub_GetRoomManager(t *testing.T) {
	hub := NewHub()

	rm := hub.GetRoomManager()

	if rm == nil {
		t.Fatal("GetRoomManager returned nil")
	}

	// Verify it's the same instance by creating a room and checking
	rm.CreateRoom("TEST01", "host-123")
	if hub.GetRoomManager().GetRoom("TEST01") == nil {
		t.Error("RoomManager should be the same instance")
	}
}

// =============================================================================
// Client Registration Tests
// =============================================================================

func TestHub_RegisterClient_NewRoom(t *testing.T) {
	// WHAT THIS TESTS: When a client registers with a room code that doesn't
	// exist yet, the hub should automatically create the room connections entry.

	hub := NewHub()
	go hub.Run() // Start the hub's event loop in a background goroutine

	// Create a test client - no real WebSocket connection needed!
	// We just need the Send channel to receive broadcast messages.
	client := &Client{
		RoomCode: "TEST01",
		PlayerID: "player-123",
		Send:     make(chan Message, 10), // Buffered to prevent blocking
	}

	// Send the client to the register channel
	// The hub.Run() goroutine will process this
	hub.Register(client)

	// Give the hub goroutine time to process the registration
	// In production code, you might use sync primitives instead
	time.Sleep(10 * time.Millisecond)

	// Verify the client is now in the room
	hub.mu.RLock()
	roomConns, exists := hub.rooms["TEST01"]
	hub.mu.RUnlock()

	if !exists {
		t.Fatal("room should exist after client registration")
	}
	if roomConns == nil {
		t.Fatal("room connections should not be nil")
	}

	roomConns.mu.RLock()
	clientCount := len(roomConns.clients)
	_, clientExists := roomConns.clients[client]
	roomConns.mu.RUnlock()

	if clientCount != 1 {
		t.Errorf("expected 1 client in room, got %d", clientCount)
	}
	if !clientExists {
		t.Error("client should be in the room's clients map")
	}
}

func TestHub_RegisterClient_ExistingRoom(t *testing.T) {
	// WHAT THIS TESTS: When multiple clients register with the same room code,
	// they should all be added to the same room connections.

	hub := NewHub()
	go hub.Run()

	client1 := &Client{RoomCode: "ROOM1", PlayerID: "p1", Send: make(chan Message, 10)}
	client2 := &Client{RoomCode: "ROOM1", PlayerID: "p2", Send: make(chan Message, 10)}
	client3 := &Client{RoomCode: "ROOM1", PlayerID: "p3", Send: make(chan Message, 10)}

	// Register all three clients to the same room
	hub.Register(client1)
	hub.Register(client2)
	hub.Register(client3)

	time.Sleep(20 * time.Millisecond)

	hub.mu.RLock()
	roomConns := hub.rooms["ROOM1"]
	hub.mu.RUnlock()

	if roomConns == nil {
		t.Fatal("room should exist")
	}

	roomConns.mu.RLock()
	clientCount := len(roomConns.clients)
	roomConns.mu.RUnlock()

	if clientCount != 3 {
		t.Errorf("expected 3 clients in room, got %d", clientCount)
	}
}

func TestHub_RegisterClient_DifferentRooms(t *testing.T) {
	// WHAT THIS TESTS: Clients in different rooms should be kept separate.

	hub := NewHub()
	go hub.Run()

	client1 := &Client{RoomCode: "ROOM1", PlayerID: "p1", Send: make(chan Message, 10)}
	client2 := &Client{RoomCode: "ROOM2", PlayerID: "p2", Send: make(chan Message, 10)}

	hub.Register(client1)
	hub.Register(client2)

	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	room1 := hub.rooms["ROOM1"]
	room2 := hub.rooms["ROOM2"]
	hub.mu.RUnlock()

	if room1 == nil || room2 == nil {
		t.Fatal("both rooms should exist")
	}

	room1.mu.RLock()
	room1Count := len(room1.clients)
	room1.mu.RUnlock()

	room2.mu.RLock()
	room2Count := len(room2.clients)
	room2.mu.RUnlock()

	if room1Count != 1 {
		t.Errorf("ROOM1 should have 1 client, got %d", room1Count)
	}
	if room2Count != 1 {
		t.Errorf("ROOM2 should have 1 client, got %d", room2Count)
	}
}

// =============================================================================
// Client Unregistration Tests
// =============================================================================

func TestHub_UnregisterClient(t *testing.T) {
	// WHAT THIS TESTS: When a client unregisters, it should be removed from
	// the room and its Send channel should be closed.

	hub := NewHub()
	go hub.Run()

	client := &Client{
		RoomCode: "TEST01",
		PlayerID: "player-123",
		Send:     make(chan Message, 10),
	}

	// Register first
	hub.Register(client)
	time.Sleep(10 * time.Millisecond)

	// Now unregister
	hub.Unregister(client)
	time.Sleep(10 * time.Millisecond)

	// Verify client's Send channel is closed
	// Sending to a closed channel would panic, but receiving returns ok=false
	select {
	case _, ok := <-client.Send:
		if ok {
			t.Error("Send channel should be closed (ok should be false)")
		}
		// ok is false, channel is closed - good!
	default:
		// Channel might be empty but open - try a different check
		// Actually, for a closed channel, we should get ok=false immediately
		t.Error("expected closed channel to return immediately")
	}
}

func TestHub_UnregisterClient_RoomCleanup(t *testing.T) {
	// WHAT THIS TESTS: When the last client leaves a room, the room should
	// be deleted from the hub's rooms map to prevent memory leaks.

	hub := NewHub()
	go hub.Run()

	client := &Client{
		RoomCode: "TEST01",
		PlayerID: "player-123",
		Send:     make(chan Message, 10),
	}

	// Register and unregister
	hub.Register(client)
	time.Sleep(10 * time.Millisecond)
	hub.Unregister(client)
	time.Sleep(10 * time.Millisecond)

	// Verify room no longer exists
	hub.mu.RLock()
	_, exists := hub.rooms["TEST01"]
	hub.mu.RUnlock()

	if exists {
		t.Error("room should be deleted after last client leaves")
	}
}

func TestHub_UnregisterClient_OtherClientsRemain(t *testing.T) {
	// WHAT THIS TESTS: When one client leaves but others remain, only that
	// client should be removed and the room should still exist.

	hub := NewHub()
	go hub.Run()

	client1 := &Client{RoomCode: "ROOM1", PlayerID: "p1", Send: make(chan Message, 10)}
	client2 := &Client{RoomCode: "ROOM1", PlayerID: "p2", Send: make(chan Message, 10)}

	// Register both
	hub.Register(client1)
	hub.Register(client2)
	time.Sleep(10 * time.Millisecond)

	// Unregister only client1
	hub.Unregister(client1)
	time.Sleep(10 * time.Millisecond)

	// Room should still exist with client2
	hub.mu.RLock()
	roomConns := hub.rooms["ROOM1"]
	hub.mu.RUnlock()

	if roomConns == nil {
		t.Fatal("room should still exist")
	}

	roomConns.mu.RLock()
	clientCount := len(roomConns.clients)
	_, client1Exists := roomConns.clients[client1]
	_, client2Exists := roomConns.clients[client2]
	roomConns.mu.RUnlock()

	if clientCount != 1 {
		t.Errorf("expected 1 client remaining, got %d", clientCount)
	}
	if client1Exists {
		t.Error("client1 should no longer be in room")
	}
	if !client2Exists {
		t.Error("client2 should still be in room")
	}
}

// =============================================================================
// Broadcast Tests
// =============================================================================

func TestHub_BroadcastToRoom_SingleClient(t *testing.T) {
	// WHAT THIS TESTS: A broadcast to a room should deliver the message
	// to the client's Send channel.

	hub := NewHub()
	go hub.Run()

	client := &Client{
		RoomCode: "TEST01",
		PlayerID: "player-123",
		Send:     make(chan Message, 10),
	}

	hub.Register(client)
	time.Sleep(10 * time.Millisecond)

	// Broadcast a message
	msg := Message{Type: "test_message", Data: nil}
	hub.BroadcastToRoom("TEST01", msg)

	// Verify client received the message
	select {
	case received := <-client.Send:
		if received.Type != "test_message" {
			t.Errorf("expected message type 'test_message', got '%s'", received.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("client did not receive broadcast within timeout")
	}
}

func TestHub_BroadcastToRoom_MultipleClients(t *testing.T) {
	// WHAT THIS TESTS: A broadcast should be delivered to ALL clients in the room.

	hub := NewHub()
	go hub.Run()

	client1 := &Client{RoomCode: "ROOM1", PlayerID: "p1", Send: make(chan Message, 10)}
	client2 := &Client{RoomCode: "ROOM1", PlayerID: "p2", Send: make(chan Message, 10)}
	client3 := &Client{RoomCode: "ROOM1", PlayerID: "p3", Send: make(chan Message, 10)}

	hub.Register(client1)
	hub.Register(client2)
	hub.Register(client3)
	time.Sleep(10 * time.Millisecond)

	// Broadcast
	msg := Message{Type: "test", Data: nil}
	hub.BroadcastToRoom("ROOM1", msg)

	// All three clients should receive the message
	clients := []*Client{client1, client2, client3}
	for i, c := range clients {
		select {
		case received := <-c.Send:
			if received.Type != "test" {
				t.Errorf("client%d: expected type 'test', got '%s'", i+1, received.Type)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("client%d: did not receive broadcast", i+1)
		}
	}
}

func TestHub_BroadcastToRoom_OnlyTargetRoom(t *testing.T) {
	// WHAT THIS TESTS: A broadcast to one room should NOT be delivered to
	// clients in other rooms.

	hub := NewHub()
	go hub.Run()

	client1 := &Client{RoomCode: "ROOM1", PlayerID: "p1", Send: make(chan Message, 10)}
	client2 := &Client{RoomCode: "ROOM2", PlayerID: "p2", Send: make(chan Message, 10)}

	hub.Register(client1)
	hub.Register(client2)
	time.Sleep(10 * time.Millisecond)

	// Broadcast only to ROOM1
	hub.BroadcastToRoom("ROOM1", Message{Type: "for_room1", Data: nil})

	// client1 should receive it
	select {
	case msg := <-client1.Send:
		if msg.Type != "for_room1" {
			t.Errorf("client1 got wrong message: %s", msg.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("client1 should have received the broadcast")
	}

	// client2 should NOT receive it
	select {
	case msg := <-client2.Send:
		t.Errorf("client2 should not receive ROOM1 broadcast, got: %s", msg.Type)
	case <-time.After(50 * time.Millisecond):
		// Good - no message received
	}
}

func TestHub_BroadcastToRoom_NonexistentRoom(t *testing.T) {
	// WHAT THIS TESTS: Broadcasting to a non-existent room should not panic.

	hub := NewHub()
	go hub.Run()

	// This should not panic
	hub.BroadcastToRoom("NONEXISTENT", Message{Type: "test", Data: nil})

	// If we get here, the test passes
}

func TestHub_BroadcastToRoom_FullSendChannel(t *testing.T) {
	// WHAT THIS TESTS: If a client's Send channel is full (slow client),
	// the broadcast should skip that client without blocking.
	// This prevents one slow client from blocking broadcasts to others.

	hub := NewHub()
	go hub.Run()

	// Create a client with a very small buffer
	slowClient := &Client{
		RoomCode: "TEST01",
		PlayerID: "slow",
		Send:     make(chan Message, 1), // Very small buffer
	}
	fastClient := &Client{
		RoomCode: "TEST01",
		PlayerID: "fast",
		Send:     make(chan Message, 10),
	}

	hub.Register(slowClient)
	hub.Register(fastClient)
	time.Sleep(10 * time.Millisecond)

	// Fill up the slow client's buffer
	slowClient.Send <- Message{Type: "filler", Data: nil}

	// Now broadcast - the slow client's buffer is full
	// This should NOT block
	done := make(chan bool)
	go func() {
		hub.BroadcastToRoom("TEST01", Message{Type: "broadcast", Data: nil})
		done <- true
	}()

	select {
	case <-done:
		// Good - broadcast completed without blocking
	case <-time.After(100 * time.Millisecond):
		t.Error("BroadcastToRoom blocked (slow client caused deadlock)")
	}

	// Fast client should still receive the message
	select {
	case msg := <-fastClient.Send:
		if msg.Type != "broadcast" {
			t.Errorf("fast client got wrong message: %s", msg.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("fast client should have received the broadcast")
	}
}

// =============================================================================
// Concurrency Tests
// =============================================================================

func TestHub_ConcurrentRegistrations(t *testing.T) {
	// WHAT THIS TESTS: Multiple goroutines registering clients simultaneously
	// should not cause race conditions or panics.

	hub := NewHub()
	go hub.Run()

	done := make(chan bool, 50)

	// Spawn 50 goroutines that each register a client
	for i := 0; i < 50; i++ {
		go func(id int) {
			client := &Client{
				RoomCode: "STRESS",
				PlayerID: string(rune('a' + id)),
				Send:     make(chan Message, 10),
			}
			hub.Register(client)
			done <- true
		}(i)
	}

	// Wait for all registrations
	for i := 0; i < 50; i++ {
		<-done
	}

	time.Sleep(50 * time.Millisecond)

	// Verify all clients are in the room
	hub.mu.RLock()
	roomConns := hub.rooms["STRESS"]
	hub.mu.RUnlock()

	if roomConns == nil {
		t.Fatal("room should exist")
	}

	roomConns.mu.RLock()
	clientCount := len(roomConns.clients)
	roomConns.mu.RUnlock()

	if clientCount != 50 {
		t.Errorf("expected 50 clients, got %d", clientCount)
	}
}

func TestHub_ConcurrentBroadcasts(t *testing.T) {
	// WHAT THIS TESTS: Multiple goroutines broadcasting simultaneously
	// should not cause race conditions.

	hub := NewHub()
	go hub.Run()

	client := &Client{
		RoomCode: "TEST",
		PlayerID: "listener",
		Send:     make(chan Message, 100), // Large buffer
	}

	hub.Register(client)
	time.Sleep(10 * time.Millisecond)

	done := make(chan bool, 10)

	// Spawn 10 goroutines that each broadcast a message
	for i := 0; i < 10; i++ {
		go func(id int) {
			hub.BroadcastToRoom("TEST", Message{Type: "concurrent", Data: nil})
			done <- true
		}(i)
	}

	// Wait for all broadcasts
	for i := 0; i < 10; i++ {
		<-done
	}

	// Count received messages
	time.Sleep(20 * time.Millisecond)
	received := 0
	for {
		select {
		case <-client.Send:
			received++
		default:
			goto countDone
		}
	}
countDone:

	if received != 10 {
		t.Errorf("expected 10 messages, got %d", received)
	}
}
