package game

import (
	"errors"
	"math/rand"
	"sync"
	"time"

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
	return code
}

func (rm *RoomManager) GetRoom(code string) (*models.Room, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, exists := rm.rooms[code]
	if !exists {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (rm *RoomManager) JoinRoom(code string, playerName string) (string, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[code]
	if !exists {
		return "", ErrRoomNotFound
	}

	playerID := uuid.New().String()
	_ = room
	_ = playerID
	_ = playerName

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
