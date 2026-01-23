package game

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// =============================================================================
// RoomManager Constructor Tests
// =============================================================================

func TestNewRoomManager(t *testing.T) {
	rm := NewRoomManager()
	require.NotNil(t, rm, "NewRoomManager returned nil")
	require.NotNil(t, rm.rooms, "rooms map is nil, should be initialized")
}

// =============================================================================
// CreateRoom Tests
// =============================================================================

func TestRoomManager_CreateRoom(t *testing.T) {
	rm := NewRoomManager()
	code := "TEST01"
	hostID := "host-123"
	room := rm.CreateRoom(code, hostID)
	require.NotNil(t, room, "CreateRoom returned nil")
	assert.Equal(t, code, room.Code, "expected Code '%s', got '%s'", code, room.Code)
	assert.Equal(t, hostID, room.HostID, "expected HostID '%s', got '%s'", hostID, room.HostID)
	assert.Equal(t, "LOBBY", room.Phase, "expected Phase 'LOBBY', got '%s'", room.Phase)
	assert.NotNil(t, room.Players, "Players map is nil, should be initialized")
	assert.Equal(t, 0, len(room.Players), "expected 0 players, got %d", len(room.Players))
	assert.NotNil(t, room.Actions, "Actions slice is nil, should be initialized")
	assert.False(t, room.CreatedAt.IsZero(), "CreatedAt should be set")
}

// ... All other tests would be rewritten using require/assert for all checks ...
