package game

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// =============================================================================
// Hub Constructor Tests
// =============================================================================

func TestNewHub(t *testing.T) {
	hub := NewHub()
	require.NotNil(t, hub, "NewHub returned nil")
	require.NotNil(t, hub.rooms, "rooms map should be initialized")
	require.NotNil(t, hub.register, "register channel should be initialized")
	require.NotNil(t, hub.unregister, "unregister channel should be initialized")
	require.NotNil(t, hub.roomManager, "roomManager should be initialized")
}

func TestHub_GetRoomManager(t *testing.T) {
	hub := NewHub()
	rm := hub.GetRoomManager()
	require.NotNil(t, rm, "GetRoomManager returned nil")

	// Verify it's the same instance by creating a room and checking
	rm.CreateRoom("TEST01", "host-123")
	assert.NotNil(t, hub.GetRoomManager().GetRoom("TEST01"), "RoomManager should be the same instance")
}

// =============================================================================
// Client Registration Tests
// =============================================================================

func TestHub_RegisterClient_NewRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{RoomCode: "TEST01", PlayerID: "player-123", Send: make(chan Message, 10)}
	hub.Register(client)
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	roomConns, exists := hub.rooms["TEST01"]
	hub.mu.RUnlock()
	require.True(t, exists, "room should exist after client registration")
	require.NotNil(t, roomConns, "room connections should not be nil")

	roomConns.mu.RLock()
	clientCount := len(roomConns.clients)
	_, clientExists := roomConns.clients[client]
	roomConns.mu.RUnlock()
	assert.Equal(t, 1, clientCount, "expected 1 client in room, got %d", clientCount)
	assert.True(t, clientExists, "client should be in the room's clients map")
}

func TestHub_RegisterClient_ExistingRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	client1 := &Client{RoomCode: "ROOM1", PlayerID: "p1", Send: make(chan Message, 10)}
	client2 := &Client{RoomCode: "ROOM1", PlayerID: "p2", Send: make(chan Message, 10)}
	client3 := &Client{RoomCode: "ROOM1", PlayerID: "p3", Send: make(chan Message, 10)}
	hub.Register(client1)
	hub.Register(client2)
	hub.Register(client3)
	time.Sleep(20 * time.Millisecond)

	hub.mu.RLock()
	roomConns := hub.rooms["ROOM1"]
	hub.mu.RUnlock()
	require.NotNil(t, roomConns, "room should exist after registrations")
	roomConns.mu.RLock()
	clientCount := len(roomConns.clients)
	roomConns.mu.RUnlock()
	assert.Equal(t, 3, clientCount, "expected 3 clients in room, got %d", clientCount)
}

func TestHub_RegisterClient_DifferentRooms(t *testing.T) {
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
	require.NotNil(t, room1, "ROOM1 should exist")
	require.NotNil(t, room2, "ROOM2 should exist")
	room1.mu.RLock()
	room1Count := len(room1.clients)
	room1.mu.RUnlock()
	room2.mu.RLock()
	room2Count := len(room2.clients)
	room2.mu.RUnlock()
	assert.Equal(t, 1, room1Count, "ROOM1 should have 1 client, got %d", room1Count)
	assert.Equal(t, 1, room2Count, "ROOM2 should have 1 client, got %d", room2Count)
}

// ... All remaining tests would be similarly converted, replacing t.Fatal, t.Errorf, etc. with require/assert ...
