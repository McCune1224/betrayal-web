package game

import (
	"sync"
	"time"
	"yourmodule/internal/models"
	"github.com/google/uuid"
)

type RoomManager struct {
	mu    sync.RWMutex
	rooms map[string]*models.Room
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*models.Room),
	}
}

func (rm *RoomManager) CreateRoom(code, hostID string) *models.Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room := &models.Room{
		Code:      code,
		HostID:    hostID,
		Phase:     "LOBBY",
		Players:   make(map[string]*models.Player),
		Actions:   []models.Action{},
		CreatedAt: time.Now(),
	}
	rm.rooms[code] = room
	return room
}

func (rm *RoomManager) GetRoom(code string) *models.Room {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.rooms[code]
}

func (rm *RoomManager) DeleteRoom(code string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, code)
}

func (rm *RoomManager) JoinRoom(code, playerName string) (*models.Room, *models.Player, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room := rm.rooms[code]
	if room == nil {
		return nil, nil, ErrRoomNotFound
	}

	playerID := uuid.New().String()
	player := &models.Player{
		ID:       playerID,
		Name:     playerName,
		IsAlive:  true,
		JoinedAt: time.Now(),
	}
	room.Players[playerID] = player

	return room, player, nil
}

func (rm *RoomManager) AdvancePhase(code string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room := rm.rooms[code]
	if room == nil {
		return ErrRoomNotFound
	}

	switch room.Phase {
	case "LOBBY":
		room.Phase = "NIGHT"
	case "NIGHT":
		room.Phase = "DAY"
	case "DAY":
		room.Phase = "NIGHT"
	}
	return nil
}
