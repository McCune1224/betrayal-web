package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomManager_CreateRoom(t *testing.T) {
	rm := NewRoomManager()

	t.Run("creates room with valid host", func(t *testing.T) {
		code := rm.CreateRoom("host-123")

		assert.NotEmpty(t, code, "room code should be generated")
		assert.Len(t, code, 6, "room code should be 6 characters")

		room, err := rm.GetRoom(code)
		require.NoError(t, err, "should be able to retrieve created room")
		assert.Equal(t, "host-123", room.HostID, "host ID should match")
		assert.Equal(t, "LOBBY", room.Phase, "phase should be LOBBY")
	})
}

func TestRoomManager_GetRoom(t *testing.T) {
	rm := NewRoomManager()
	code := rm.CreateRoom("host-456")

	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "existing room",
			code:    code,
			wantErr: false,
		},
		{
			name:    "non-existing room",
			code:    "FAKE01",
			wantErr: true,
		},
		{
			name:    "empty code",
			code:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room, err := rm.GetRoom(tt.code)

			if tt.wantErr {
				assert.Error(t, err, "should return error for non-existing room")
				assert.Nil(t, room, "room should be nil when not found")
			} else {
				require.NoError(t, err, "should not return error for existing room")
				assert.NotNil(t, room, "room should not be nil")
			}
		})
	}
}

func TestRoomManager_JoinRoom(t *testing.T) {
	rm := NewRoomManager()
	code := rm.CreateRoom("host-789")

	t.Run("joins existing room", func(t *testing.T) {
		playerID, err := rm.JoinRoom(code, "Alice")

		require.NoError(t, err, "should join existing room without error")
		assert.NotEmpty(t, playerID, "player ID should be generated")
	})

	t.Run("fails to join non-existing room", func(t *testing.T) {
		playerID, err := rm.JoinRoom("FAKE99", "Bob")

		assert.Error(t, err, "should return error for non-existing room")
		assert.Empty(t, playerID, "player ID should be empty on error")
		assert.Equal(t, ErrRoomNotFound, err, "error should be ErrRoomNotFound")
	})
}

func TestGenerateRoomCode(t *testing.T) {
	t.Run("generates unique codes", func(t *testing.T) {
		codes := make(map[string]bool)

		for i := 0; i < 100; i++ {
			code := generateRoomCode()

			assert.Len(t, code, 6, "code should be 6 characters")
			assert.NotContains(t, codes, code, "code should be unique")
			codes[code] = true
		}
	})

	t.Run("contains only valid characters", func(t *testing.T) {
		validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

		for i := 0; i < 50; i++ {
			code := generateRoomCode()

			for _, char := range code {
				assert.Contains(t, validChars, string(char),
					"code should only contain uppercase letters and digits")
			}
		}
	})
}
