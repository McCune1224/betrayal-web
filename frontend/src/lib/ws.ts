// WebSocket client for real-time communication
import type { 
	PlayerJoinedData, 
	PlayerLeftData, 
	PhaseChangedData, 
	ActionData, 
	GameEndedData, 
	ErrorData,
	WSMessage 
} from './types';
import { 
	setConnectionStatus, 
	updateLastPing,
	addPlayerToRoom,
	removePlayerFromRoom,
	setRoomPhase,
	addSystemMessage
} from './stores.svelte';

const WS_BASE = import.meta.env.VITE_WS_BASE || 'ws://localhost:8080';

let ws: WebSocket | null = null;
let reconnectAttempts = 0;
let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_DELAY = 3000;

/**
 * Connect to WebSocket server
 * @param roomCode - Room code to join
 * @param playerId - Player ID
 * @param playerName - Player display name
 * @returns Cleanup function to disconnect
 */
export function connect(roomCode: string, playerId: string, playerName: string): () => void {
	// Disconnect existing connection if any
	disconnect();
	
	setConnectionStatus('connecting');
	
	const wsUrl = `${WS_BASE}/ws?room=${encodeURIComponent(roomCode)}&player=${encodeURIComponent(playerId)}&name=${encodeURIComponent(playerName)}`;
	
	try {
		ws = new WebSocket(wsUrl);
		
		ws.onopen = () => {
			console.log('WebSocket connected');
			setConnectionStatus('connected');
			reconnectAttempts = 0;
			addSystemMessage('Connected to room', 'success');
		};
		
		ws.onmessage = (event: MessageEvent) => {
			handleMessage(event.data);
		};
		
		ws.onclose = (event: CloseEvent) => {
			console.log('WebSocket closed:', event.code, event.reason);
			ws = null;
			
			if (event.code === 1000 || event.code === 1001) {
				// Normal closure
				setConnectionStatus('disconnected');
				addSystemMessage('Disconnected from room', 'info');
			} else {
				// Abnormal closure - try to reconnect
				setConnectionStatus('disconnected');
				addSystemMessage('Connection lost. Trying to reconnect...', 'warning');
				attemptReconnect(roomCode, playerId, playerName);
			}
		};
		
		ws.onerror = (error: Event) => {
			console.error('WebSocket error:', error);
			setConnectionStatus('error', 'Connection error');
			addSystemMessage('Connection error occurred', 'error');
		};
		
	} catch (error) {
		console.error('Failed to connect WebSocket:', error);
		const errorMessage = error instanceof Error ? error.message : 'Unknown error';
		setConnectionStatus('error', errorMessage);
	}
	
	// Return cleanup function
	return () => disconnect();
}

/**
 * Disconnect from WebSocket server
 */
export function disconnect(): void {
	if (reconnectTimeout) {
		clearTimeout(reconnectTimeout);
		reconnectTimeout = null;
	}
	
	if (ws) {
		// Only close if still open or connecting
		if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
			ws.close(1000, 'Client disconnected');
		}
		ws = null;
	}
	
	setConnectionStatus('disconnected');
}

/**
 * Send a message to the server
 * @param message - Message object with type and data
 */
export function send(message: WSMessage): boolean {
	if (!ws || ws.readyState !== WebSocket.OPEN) {
		console.error('WebSocket is not connected');
		return false;
	}
	
	try {
		ws.send(JSON.stringify(message));
		return true;
	} catch (error) {
		console.error('Failed to send message:', error);
		return false;
	}
}

/**
 * Send a chat message
 * @param text - Message text
 */
export function sendChatMessage(text: string): boolean {
	return send({
		type: 'chat_message',
		data: { text }
	});
}

/**
 * Send a ping (keepalive)
 */
export function sendPing(): boolean {
	return send({ type: 'ping' });
}

/**
 * Handle incoming WebSocket messages
 * @param data - Raw message data
 */
function handleMessage(data: string): void {
	try {
		const message: WSMessage = JSON.parse(data);
		console.log('Received message:', message);
		
		switch (message.type) {
			case 'player_joined':
				handlePlayerJoined(message.data as PlayerJoinedData);
				break;
				
			case 'player_rejoined':
				handlePlayerRejoined(message.data as PlayerJoinedData);
				break;
				
			case 'player_left':
				handlePlayerLeft(message.data as PlayerLeftData);
				break;
				
			case 'phase_changed':
				handlePhaseChanged(message.data as PhaseChangedData);
				break;
				
			case 'action_submitted':
				handleActionSubmitted(message.data as ActionData);
				break;
				
			case 'action_deleted':
				handleActionDeleted(message.data);
				break;
				
			case 'roles_assigned':
				handleRolesAssigned(message.data);
				break;
				
			case 'game_ended':
				handleGameEnded(message.data as GameEndedData);
				break;
				
			case 'pong':
				updateLastPing();
				break;
				
			case 'error':
				handleError(message.data as ErrorData);
				break;
				
			default:
				console.log('Unknown message type:', message.type);
		}
	} catch (error) {
		console.error('Failed to parse message:', error);
	}
}

function handlePlayerJoined(data: PlayerJoinedData): void {
	addPlayerToRoom({
		id: data.playerId,
		name: data.playerName,
		isHost: data.isHost
	});
	addSystemMessage(`${data.playerName} joined the room`, 'info');
}

function handlePlayerRejoined(data: PlayerJoinedData): void {
	addPlayerToRoom({
		id: data.playerId,
		name: data.playerName,
		isHost: data.isHost
	});
	addSystemMessage(`${data.playerName} reconnected`, 'success');
}

function handlePlayerLeft(data: PlayerLeftData): void {
	removePlayerFromRoom(data.playerId);
	addSystemMessage(`${data.playerName || 'A player'} left the room`, 'info');
}

function handlePhaseChanged(data: PhaseChangedData): void {
	setRoomPhase(data.phase);
	addSystemMessage(`Phase changed to ${data.phase}`, 'info');
}

function handleActionSubmitted(data: ActionData): void {
	addSystemMessage(`${data.playerName} submitted an action`, 'info');
}

function handleActionDeleted(_data: unknown): void {
	addSystemMessage('An action was deleted by the host', 'info');
}

function handleRolesAssigned(_data: unknown): void {
	addSystemMessage('Roles have been assigned! Check your role.', 'success');
	// Could update player role here if included in data
}

function handleGameEnded(data: GameEndedData): void {
	addSystemMessage(`Game Over! ${data.winner} wins!`, 'success');
}

function handleError(data: ErrorData): void {
	console.error('Server error:', data);
	addSystemMessage(`Error: ${data.message}`, 'error');
}

/**
 * Attempt to reconnect with exponential backoff
 * @param roomCode - Room code
 * @param playerId - Player ID
 * @param playerName - Player name
 */
function attemptReconnect(roomCode: string, playerId: string, playerName: string): void {
	if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
		addSystemMessage('Failed to reconnect. Please refresh the page.', 'error');
		setConnectionStatus('error', 'Max reconnection attempts reached');
		return;
	}
	
	reconnectAttempts++;
	const delay = RECONNECT_DELAY * Math.pow(2, reconnectAttempts - 1);
	
	console.log(`Reconnecting in ${delay}ms (attempt ${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS})`);
	
	reconnectTimeout = setTimeout(() => {
		connect(roomCode, playerId, playerName);
	}, delay);
}

type ConnectionState = 'connecting' | 'open' | 'closing' | 'closed' | 'unknown';

/**
 * Get current WebSocket state
 * @returns Current connection state
 */
export function getConnectionState(): ConnectionState {
	if (!ws) return 'closed';
	
	switch (ws.readyState) {
		case WebSocket.CONNECTING: return 'connecting';
		case WebSocket.OPEN: return 'open';
		case WebSocket.CLOSING: return 'closing';
		case WebSocket.CLOSED: return 'closed';
		default: return 'unknown';
	}
}
