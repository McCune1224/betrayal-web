package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	HostID    string    `json:"hostId"`
	Phase     string    `json:"phase"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Player struct {
	ID       uuid.UUID `json:"id"`
	RoomID   uuid.UUID `json:"roomId"`
	Name     string    `json:"name"`
	RoleID   *int      `json:"roleId,omitempty"`
	IsAlive  bool      `json:"isAlive"`
	JoinedAt time.Time `json:"joinedAt"`
}

type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Team        string    `json:"team"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Action struct {
	ID         uuid.UUID `json:"id"`
	RoomID     uuid.UUID `json:"roomId"`
	PlayerID   uuid.UUID `json:"playerId"`
	ActionType string    `json:"actionType"`
	TargetID   uuid.UUID `json:"targetId,omitempty"`
	Phase      string    `json:"phase"`
	CreatedAt  time.Time `json:"createdAt"`
}
