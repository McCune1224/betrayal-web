package game

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"backend/internal/logging"
	"backend/internal/models"
	"github.com/google/uuid"
)

var ErrRoomNotFound = errors.New("room not found")
var ErrPlayerNotFound = errors.New("player not found")

type RoomManager struct {
	mu    sync.RWMutex
	rooms map[string]*models.Room
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*models.Room),
	}
}

func (rm *RoomManager) CreateRoom(hostID string) string {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	code := generateRoomCode()
	room := &models.Room{
		ID:        uuid.New(),
		Code:      code,
		HostID:    hostID,
		Phase:     "LOBBY",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rm.rooms[code] = room

	logger := logging.RoomLogger(code)
	logger.Info("room_created",
		"host_id", hostID,
		"total_rooms", len(rm.rooms),
	)

	return code
}

func (rm *RoomManager) GetRoom(code string) (*models.Room, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, exists := rm.rooms[code]
	if !exists {
		logger := logging.RoomLogger(code)
		logger.Debug("room_not_found")
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (rm *RoomManager) JoinRoom(code string, playerName string) (string, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[code]
	if !exists {
		logger := logging.RoomLogger(code)
		logger.Warn("join_attempt_to_nonexistent_room",
			"player_name", playerName,
		)
		return "", ErrRoomNotFound
	}

	playerID := uuid.New().String()
	_ = room
	_ = playerID
	_ = playerName

	logger := logging.RoomLogger(code)
	logger.Info("player_joined",
		"player_id", playerID,
		"player_name", playerName,
		"room_phase", room.Phase,
	)

	return playerID, nil
}

func generateRoomCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	rand.Seed(time.Now().UnixNano())
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
