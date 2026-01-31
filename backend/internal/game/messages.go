package game

import "time"

// Message types for WebSocket communication
const (
	// Connection events
	MsgTypePlayerJoined   = "player_joined"
	MsgTypePlayerLeft     = "player_left"
	MsgTypePlayerRejoined = "player_rejoined"

	// Game phase events
	MsgTypePhaseChanged = "phase_changed"
	MsgTypeGameStarted  = "game_started"
	MsgTypeGameEnded    = "game_ended"

	// Action events
	MsgTypeActionSubmitted = "action_submitted"
	MsgTypeActionDeleted   = "action_deleted"
	MsgTypeActionsCleared  = "actions_cleared"

	// Role events
	MsgTypeRolesAssigned = "roles_assigned"
	MsgTypeRoleRevealed  = "role_revealed"

	// Host events
	MsgTypeHostChanged  = "host_changed"
	MsgTypePlayerKicked = "player_kicked"

	// Error events
	MsgTypeError         = "error"
	MsgTypeSystemMessage = "system_message"
)

// WSMessage is the envelope for all WebSocket messages
type WSMessage struct {
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// PlayerJoinedData is sent when a player joins
type PlayerJoinedData struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
	IsHost     bool   `json:"isHost"`
}

// PlayerLeftData is sent when a player leaves
type PlayerLeftData struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
}

// PhaseChangedData is sent when game phase changes
type PhaseChangedData struct {
	PreviousPhase string `json:"previousPhase"`
	CurrentPhase  string `json:"currentPhase"`
}

// ActionSubmittedData is sent when a player submits an action
type ActionSubmittedData struct {
	ActionID   string `json:"actionId"`
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
	ActionType string `json:"actionType"`
	TargetID   string `json:"targetId,omitempty"`
}

// ActionDeletedData is sent when host deletes an action
type ActionDeletedData struct {
	ActionID string `json:"actionId"`
}

// RolesAssignedData is sent when game starts and roles are assigned
type RolesAssignedData struct {
	Players []PlayerRoleInfo `json:"players"`
}

// PlayerRoleInfo contains role info for a player
type PlayerRoleInfo struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
	RoleID     int    `json:"roleId"`
	RoleName   string `json:"roleName"`
	Team       string `json:"team"`
}

// ErrorData is sent when an error occurs
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SystemMessageData is for general system notifications
type SystemMessageData struct {
	Message string `json:"message"`
	Level   string `json:"level"` // info, warning, error
}

// Incoming message types (from clients)

// JoinRoomMessage is sent by client to join a room
type JoinRoomMessage struct {
	RoomCode   string `json:"roomCode"`
	PlayerName string `json:"playerName"`
}

// SubmitActionMessage is sent by client to submit an action
type SubmitActionMessage struct {
	ActionType string `json:"actionType"`
	TargetID   string `json:"targetId,omitempty"`
}

// HostCommandMessage is sent by host for game management
type HostCommandMessage struct {
	Command string      `json:"command"` // start_game, advance_phase, kick_player, delete_action
	Data    interface{} `json:"data,omitempty"`
}

// NewWSMessage creates a new WebSocket message with current timestamp
func NewWSMessage(msgType string, data interface{}) WSMessage {
	return WSMessage{
		Type:      msgType,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
}
