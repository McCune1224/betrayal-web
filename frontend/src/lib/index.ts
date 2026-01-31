// place files you want to import through the `$lib` alias in this folder.

// Re-export all types
export type {
	Player,
	RoomPlayer,
	Room,
	Connection,
	Message,
	WSMessage,
	PlayerJoinedData,
	PlayerLeftData,
	PhaseChangedData,
	ActionData,
	GameEndedData,
	ErrorData,
	CreateRoomResponse,
	JoinRoomResponse,
	HealthCheckResponse
} from './types';

// Re-export stores and functions
export {
	player,
	room,
	connection,
	messages,
	log,
	getIsHost,
	setPlayer,
	clearPlayer,
	loadPlayerFromStorage,
	setRoom,
	clearRoom,
	loadRoomFromStorage,
	addPlayerToRoom,
	removePlayerFromRoom,
	setConnectionStatus,
	updateLastPing,
	addMessage,
	addSystemMessage,
	addLogEntry,
	clearMessages,
	clearLog,
	setRoomPhase
} from './stores.svelte';

// Re-export API functions
export { createRoom, joinRoom, healthCheck } from './api';

// Re-export WebSocket functions
export { connect, disconnect, send, sendChatMessage, sendPing, getConnectionState } from './ws';
