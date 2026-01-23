package game

import (
	"testing"
)

// =============================================================================
// RoomManager Constructor Tests
// =============================================================================

func TestNewRoomManager(t *testing.T) {
	// NewRoomManager should create an empty RoomManager with initialized map
	rm := NewRoomManager()

	if rm == nil {
		t.Fatal("NewRoomManager returned nil")
	}

	// The rooms map should exist but be empty
	if rm.rooms == nil {
		t.Fatal("rooms map is nil, should be initialized")
	}
}

// =============================================================================
// CreateRoom Tests
// =============================================================================

func TestRoomManager_CreateRoom(t *testing.T) {
	rm := NewRoomManager()

	code := "TEST01"
	hostID := "host-123"

	room := rm.CreateRoom(code, hostID)

	// Verify room was created with correct fields
	if room == nil {
		t.Fatal("CreateRoom returned nil")
	}
	if room.Code != code {
		t.Errorf("expected Code '%s', got '%s'", code, room.Code)
	}
	if room.HostID != hostID {
		t.Errorf("expected HostID '%s', got '%s'", hostID, room.HostID)
	}
	if room.Phase != "LOBBY" {
		t.Errorf("expected Phase 'LOBBY', got '%s'", room.Phase)
	}
	if room.Players == nil {
		t.Error("Players map is nil, should be initialized")
	}
	if len(room.Players) != 0 {
		t.Errorf("expected 0 players, got %d", len(room.Players))
	}
	if room.Actions == nil {
		t.Error("Actions slice is nil, should be initialized")
	}
	if room.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestRoomManager_CreateRoom_MultipleDifferentCodes(t *testing.T) {
	rm := NewRoomManager()

	// Create multiple rooms with different codes
	room1 := rm.CreateRoom("ROOM01", "host-1")
	room2 := rm.CreateRoom("ROOM02", "host-2")
	room3 := rm.CreateRoom("ROOM03", "host-3")

	// All rooms should exist independently
	if rm.GetRoom("ROOM01") != room1 {
		t.Error("ROOM01 not found or doesn't match")
	}
	if rm.GetRoom("ROOM02") != room2 {
		t.Error("ROOM02 not found or doesn't match")
	}
	if rm.GetRoom("ROOM03") != room3 {
		t.Error("ROOM03 not found or doesn't match")
	}
}

// =============================================================================
// GetRoom Tests
// =============================================================================

func TestRoomManager_GetRoom_Exists(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	room := rm.GetRoom("TEST01")

	if room == nil {
		t.Fatal("GetRoom returned nil for existing room")
	}
	if room.Code != "TEST01" {
		t.Errorf("expected Code 'TEST01', got '%s'", room.Code)
	}
}

func TestRoomManager_GetRoom_NotExists(t *testing.T) {
	rm := NewRoomManager()

	room := rm.GetRoom("NONEXISTENT")

	if room != nil {
		t.Error("GetRoom should return nil for non-existent room")
	}
}

// =============================================================================
// DeleteRoom Tests
// =============================================================================

func TestRoomManager_DeleteRoom(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	// Verify room exists before delete
	if rm.GetRoom("TEST01") == nil {
		t.Fatal("room should exist before delete")
	}

	rm.DeleteRoom("TEST01")

	// Verify room no longer exists
	if rm.GetRoom("TEST01") != nil {
		t.Error("room should be nil after delete")
	}
}

func TestRoomManager_DeleteRoom_NonExistent(t *testing.T) {
	rm := NewRoomManager()

	// Deleting non-existent room should not panic
	rm.DeleteRoom("NONEXISTENT")

	// If we get here without panic, test passes
}

// =============================================================================
// JoinRoom Tests
// =============================================================================

func TestRoomManager_JoinRoom_Success(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	room, player, err := rm.JoinRoom("TEST01", "Alice")

	if err != nil {
		t.Fatalf("JoinRoom returned error: %v", err)
	}
	if room == nil {
		t.Fatal("JoinRoom returned nil room")
	}
	if player == nil {
		t.Fatal("JoinRoom returned nil player")
	}

	// Verify player fields
	if player.ID == "" {
		t.Error("player ID should not be empty")
	}
	if player.Name != "Alice" {
		t.Errorf("expected player name 'Alice', got '%s'", player.Name)
	}
	if !player.IsAlive {
		t.Error("player should be alive on join")
	}
	if player.JoinedAt.IsZero() {
		t.Error("JoinedAt should be set")
	}

	// Verify player is in room
	if _, exists := room.Players[player.ID]; !exists {
		t.Error("player should be in room.Players map")
	}
}

func TestRoomManager_JoinRoom_RoomNotFound(t *testing.T) {
	rm := NewRoomManager()

	room, player, err := rm.JoinRoom("NONEXISTENT", "Alice")

	if err != ErrRoomNotFound {
		t.Errorf("expected ErrRoomNotFound, got %v", err)
	}
	if room != nil {
		t.Error("room should be nil when error occurs")
	}
	if player != nil {
		t.Error("player should be nil when error occurs")
	}
}

func TestRoomManager_JoinRoom_MultiplePlayers(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	// Join multiple players
	_, player1, _ := rm.JoinRoom("TEST01", "Alice")
	_, player2, _ := rm.JoinRoom("TEST01", "Bob")
	room, player3, _ := rm.JoinRoom("TEST01", "Charlie")

	// Verify all players are in room
	if len(room.Players) != 3 {
		t.Errorf("expected 3 players, got %d", len(room.Players))
	}

	// Verify each player has unique ID
	ids := map[string]bool{
		player1.ID: true,
		player2.ID: true,
		player3.ID: true,
	}
	if len(ids) != 3 {
		t.Error("player IDs should be unique")
	}
}

// =============================================================================
// AdvancePhase Tests
// =============================================================================

func TestRoomManager_AdvancePhase_LobbyToNight(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	err := rm.AdvancePhase("TEST01")

	if err != nil {
		t.Fatalf("AdvancePhase returned error: %v", err)
	}

	room := rm.GetRoom("TEST01")
	if room.Phase != "NIGHT" {
		t.Errorf("expected Phase 'NIGHT', got '%s'", room.Phase)
	}
}

func TestRoomManager_AdvancePhase_NightToDay(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	// Advance to NIGHT first
	rm.AdvancePhase("TEST01")

	// Then advance to DAY
	err := rm.AdvancePhase("TEST01")

	if err != nil {
		t.Fatalf("AdvancePhase returned error: %v", err)
	}

	room := rm.GetRoom("TEST01")
	if room.Phase != "DAY" {
		t.Errorf("expected Phase 'DAY', got '%s'", room.Phase)
	}
}

func TestRoomManager_AdvancePhase_DayToNight(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	// Advance through LOBBY -> NIGHT -> DAY
	rm.AdvancePhase("TEST01") // NIGHT
	rm.AdvancePhase("TEST01") // DAY

	// Then advance back to NIGHT
	err := rm.AdvancePhase("TEST01")

	if err != nil {
		t.Fatalf("AdvancePhase returned error: %v", err)
	}

	room := rm.GetRoom("TEST01")
	if room.Phase != "NIGHT" {
		t.Errorf("expected Phase 'NIGHT', got '%s'", room.Phase)
	}
}

func TestRoomManager_AdvancePhase_FullCycle(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	// Test a full cycle: LOBBY -> NIGHT -> DAY -> NIGHT -> DAY
	expectedPhases := []string{"NIGHT", "DAY", "NIGHT", "DAY", "NIGHT"}

	for i, expected := range expectedPhases {
		rm.AdvancePhase("TEST01")
		room := rm.GetRoom("TEST01")
		if room.Phase != expected {
			t.Errorf("step %d: expected Phase '%s', got '%s'", i+1, expected, room.Phase)
		}
	}
}

func TestRoomManager_AdvancePhase_RoomNotFound(t *testing.T) {
	rm := NewRoomManager()

	err := rm.AdvancePhase("NONEXISTENT")

	if err != ErrRoomNotFound {
		t.Errorf("expected ErrRoomNotFound, got %v", err)
	}
}

// =============================================================================
// Concurrency Tests
// =============================================================================

func TestRoomManager_ConcurrentJoins(t *testing.T) {
	// This test verifies that multiple goroutines can safely join the same room
	// concurrently without causing data races or panics.

	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	// Use a channel to synchronize goroutines
	done := make(chan bool, 10)

	// Spawn 10 goroutines that all try to join at the same time
	for i := 0; i < 10; i++ {
		go func(playerNum int) {
			_, _, err := rm.JoinRoom("TEST01", "Player")
			if err != nil {
				t.Errorf("goroutine %d: JoinRoom failed: %v", playerNum, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all 10 players were added
	room := rm.GetRoom("TEST01")
	if len(room.Players) != 10 {
		t.Errorf("expected 10 players, got %d", len(room.Players))
	}
}

func TestRoomManager_ConcurrentPhaseAdvances(t *testing.T) {
	// This test verifies that concurrent phase advances don't cause panics.
	// Note: The final phase may vary depending on timing, but no race should occur.

	rm := NewRoomManager()
	rm.CreateRoom("TEST01", "host-123")

	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func() {
			rm.AdvancePhase("TEST01")
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}

	// If we get here without panic or race detector errors, the test passes.
	// The actual phase depends on timing, but it should be a valid phase.
	room := rm.GetRoom("TEST01")
	validPhases := map[string]bool{"LOBBY": true, "NIGHT": true, "DAY": true}
	if !validPhases[room.Phase] {
		t.Errorf("unexpected phase: %s", room.Phase)
	}
}
