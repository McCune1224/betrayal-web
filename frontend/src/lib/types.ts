// TypeScript type definitions for the application

export interface Player {
	id: string | null;
	name: string;
	isHost: boolean;
	joinedAt: string | null;
}

export interface RoomPlayer {
	id: string;
	name: string;
	isHost: boolean;
}

export interface Room {
	code: string;
	phase: string;
	players: RoomPlayer[];
	hostId: string | null;
}

export interface Connection {
	status: 'idle' | 'connecting' | 'connected' | 'disconnected' | 'error';
	error: string | null;
	lastPing: string | null;
}

export interface Message {
	id: string;
	type: 'chat' | 'system';
	subtype?: string;
	sender?: string;
	senderId?: string;
	text?: string;
	timestamp: string;
}

// WebSocket message types
export interface WSMessage {
	type: string;
	data?: any;
}

export interface PlayerJoinedData {
	playerId: string;
	playerName: string;
	isHost: boolean;
}

export interface PlayerLeftData {
	playerId: string;
	playerName?: string;
}

export interface PhaseChangedData {
	phase: string;
}

export interface ActionData {
	playerName: string;
}

export interface GameEndedData {
	winner: string;
}

export interface ErrorData {
	message: string;
}

// API response types
export interface CreateRoomResponse {
	roomCode: string;
	hostId: string;
}

export interface JoinRoomResponse {
	playerId: string;
	phase: string;
}

export interface HealthCheckResponse {
	status: string;
}
