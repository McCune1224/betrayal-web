// Svelte 5 runes-based stores for game state management
import type { Player, RoomPlayer, Room, Connection, Message } from './types';

// Player store - current user's information
export const player = $state<Player>({
	id: null,
	name: '',
	isHost: false,
	joinedAt: null
});

// Room store - room state and player list
export const room = $state<Room>({
	code: '',
	phase: 'LOBBY',
	players: [],
	hostId: null
});

// Connection store - WebSocket connection status
export const connection = $state<Connection>({
	status: 'idle',
	error: null,
	lastPing: null
});

// Messages store - chat and game messages
export const messages = $state<Message[]>([]);

// Event log store - system events
export const log = $state<Array<{id: string; type?: string; text?: string; timestamp: string}>>([]);

// Derived value for checking if current player is host
// Note: $derived cannot be exported directly, so we use a getter function
export function getIsHost(): boolean {
	return Boolean(player.id && room.hostId && player.id === room.hostId);
}

// Helper functions to update stores

export function setPlayer(id: string, name: string, isHostPlayer: boolean = false): void {
	player.id = id;
	player.name = name;
	player.isHost = isHostPlayer;
	player.joinedAt = new Date().toISOString();
	
	// Persist player ID for reconnection
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem('playerId', id);
		localStorage.setItem('playerName', name);
	}
}

export function clearPlayer(): void {
	player.id = null;
	player.name = '';
	player.isHost = false;
	player.joinedAt = null;
	
	if (typeof localStorage !== 'undefined') {
		localStorage.removeItem('playerId');
		localStorage.removeItem('playerName');
	}
}

export function loadPlayerFromStorage(): boolean {
	if (typeof localStorage !== 'undefined') {
		const id = localStorage.getItem('playerId');
		const name = localStorage.getItem('playerName');
		if (id && name) {
			player.id = id;
			player.name = name;
			return true;
		}
	}
	return false;
}

export function setRoom(code: string, phase: string = 'LOBBY', hostId: string | null = null): void {
	room.code = code;
	room.phase = phase;
	room.hostId = hostId;
	
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem('roomCode', code);
	}
}

export function clearRoom(): void {
	room.code = '';
	room.phase = 'LOBBY';
	room.players = [];
	room.hostId = null;
	
	if (typeof localStorage !== 'undefined') {
		localStorage.removeItem('roomCode');
	}
}

export function loadRoomFromStorage(): boolean {
	if (typeof localStorage !== 'undefined') {
		const code = localStorage.getItem('roomCode');
		if (code) {
			room.code = code;
			return true;
		}
	}
	return false;
}

export function addPlayerToRoom(playerData: RoomPlayer): void {
	// Check if player already exists
	const existingIndex = room.players.findIndex(p => p.id === playerData.id);
	if (existingIndex >= 0) {
		// Update existing player
		room.players[existingIndex] = { ...room.players[existingIndex], ...playerData };
	} else {
		// Add new player
		room.players.push(playerData);
	}
}

export function removePlayerFromRoom(playerId: string): void {
	room.players = room.players.filter(p => p.id !== playerId);
}

export function setConnectionStatus(status: Connection['status'], error: string | null = null): void {
	connection.status = status;
	connection.error = error;
}

export function updateLastPing(): void {
	connection.lastPing = new Date().toISOString();
}

interface ChatMessage {
	type: 'chat';
	sender: string;
	senderId: string;
	text: string;
}

export function addMessage(message: ChatMessage): void {
	messages.push({
		...message,
		id: crypto.randomUUID(),
		timestamp: new Date().toISOString()
	});
}

export function addSystemMessage(text: string, type: string = 'info'): void {
	messages.push({
		id: crypto.randomUUID(),
		type: 'system',
		subtype: type,
		text,
		timestamp: new Date().toISOString()
	});
}

export function addLogEntry(entry: Record<string, any>): void {
	log.push({
		...entry,
		id: crypto.randomUUID(),
		timestamp: new Date().toISOString()
	});
}

export function clearMessages(): void {
	messages.length = 0;
}

export function clearLog(): void {
	log.length = 0;
}

export function setRoomPhase(phase: string): void {
	room.phase = phase;
}
