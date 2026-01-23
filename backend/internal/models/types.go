package models

import "time"

type Room struct {
	Code      string
	HostID    string
	Phase     string
	Players   map[string]*Player
	Actions   []Action
	CreatedAt time.Time
}

type Player struct {
	ID       string
	Name     string
	RoleID   int
	IsAlive  bool
	JoinedAt time.Time
}

type Action struct {
	ID        string
	PlayerID  string
	Type      string
	TargetID  string
	Phase     string
	Timestamp time.Time
}

type Role struct {
	ID       int
	Name     string
	Alignment string
	Perks    string
}
