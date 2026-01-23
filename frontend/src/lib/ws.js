// =============================================================================
// WebSocket Client
// =============================================================================
//
// This module manages the WebSocket connection to the backend.
// It handles connecting, disconnecting, sending messages, and updating
// Svelte stores based on incoming messages.
//
// KEY CONCEPTS:
// 1. WebSocket is a persistent, bidirectional connection (unlike HTTP)
// 2. Messages are sent/received as JSON strings
// 3. We update Svelte stores when messages arrive, triggering UI updates
// 4. Automatic reconnection is attempted on disconnect
//
// =============================================================================

import { room, player, actions, connectionStatus, addLogEvent } from './stores.js';

/**
 * Base URL for the WebSocket server.
 * In development: ws://localhost:8080
 * In production: Set via VITE_WS_BASE environment variable
 */
const WS_BASE = import.meta.env.VITE_WS_BASE || 'ws://localhost:8080';

// =============================================================================
// Module State
// =============================================================================

/** @type {WebSocket | null} */
let socket = null;

/** @type {string | null} */
let currentRoomCode = null;

/** @type {string | null} */
let currentPlayerId = null;

/** @type {number | null} */
let reconnectTimer = null;

/** Number of reconnection attempts */
let reconnectAttempts = 0;

/** Maximum reconnection attempts before giving up */
const MAX_RECONNECT_ATTEMPTS = 5;

/** Delay between reconnection attempts (increases with each attempt) */
const RECONNECT_DELAY_MS = 2000;

// =============================================================================
// Public API
// =============================================================================

/**
 * Connects to the WebSocket server for a specific room.
 *
 * @param {string} roomCode - The room code to join
 * @param {string} playerId - The player's ID
 * @returns {void}
 *
 * @example
 * connectWS("abc123", "player-xyz");
 */
export function connectWS(roomCode, playerId) {
	// Close existing connection if any
	if (socket) {
		closeWS();
	}

	currentRoomCode = roomCode;
	currentPlayerId = playerId;
	reconnectAttempts = 0;

	doConnect();
}

/**
 * Closes the WebSocket connection.
 * Call this when leaving a room or unmounting the component.
 */
export function closeWS() {
	// Clear any pending reconnection
	if (reconnectTimer) {
		clearTimeout(reconnectTimer);
		reconnectTimer = null;
	}

	if (socket) {
		// Remove event handlers to prevent reconnection
		socket.onclose = null;
		socket.onerror = null;
		socket.close();
		socket = null;
	}

	currentRoomCode = null;
	currentPlayerId = null;
	connectionStatus.set('disconnected');
}

/**
 * Sends a message to the server via WebSocket.
 *
 * @param {Object} message - The message to send
 * @param {string} message.type - Message type (e.g., "chat", "action")
 * @param {*} [message.data] - Message payload
 * @returns {boolean} True if sent, false if not connected
 *
 * @example
 * sendMessage({ type: "chat", data: "Hello everyone!" });
 * sendMessage({ type: "action", data: { type: "vote", targetId: "player-123" } });
 */
export function sendMessage(message) {
	if (!socket || socket.readyState !== WebSocket.OPEN) {
		console.warn('WebSocket not connected, cannot send message');
		return false;
	}

	socket.send(JSON.stringify(message));
	return true;
}

/**
 * Sends a chat message to the room.
 *
 * @param {string} text - The chat message text
 * @returns {boolean} True if sent
 */
export function sendChat(text) {
	return sendMessage({ type: 'chat', data: text });
}

/**
 * Requests the host to advance the game phase.
 * Only works if the current player is the host.
 *
 * @returns {boolean} True if sent
 */
export function sendAdvancePhase() {
	return sendMessage({ type: 'advance_phase', data: null });
}

/**
 * Requests the host to start the game (from LOBBY to NIGHT).
 * Only works if the current player is the host and there are 3+ players.
 *
 * @returns {boolean} True if sent
 */
export function sendStartGame() {
	return sendMessage({ type: 'start_game', data: null });
}

/**
 * Submits a game action (vote, kill, investigate, protect).
 *
 * @param {string} actionType - The action type (e.g., "vote", "kill")
 * @param {string} [targetId] - The target player ID (if applicable)
 * @returns {boolean} True if sent
 */
export function sendAction(actionType, targetId = undefined) {
	return sendMessage({
		type: 'submit_action',
		data: { type: actionType, target_id: targetId || null }
	});
}

/**
 * Requests deletion of an action (host only).
 *
 * @param {string} actionId - The action ID to delete
 * @returns {boolean} True if sent
 */
export function sendDeleteAction(actionId) {
	return sendMessage({
		type: 'delete_action',
		data: { action_id: actionId }
	});
}

// =============================================================================
// Internal Functions
// =============================================================================

/**
 * Creates the WebSocket connection.
 * @private
 */
function doConnect() {
	if (!currentRoomCode || !currentPlayerId) {
		console.error('Cannot connect: missing room code or player ID');
		return;
	}

	connectionStatus.set('connecting');

	const wsUrl = `${WS_BASE}/ws?room=${currentRoomCode}&player=${currentPlayerId}`;
	socket = new WebSocket(wsUrl);

	// --- Event: Connection Opened ---
	socket.onopen = () => {
		console.log('WebSocket connected');
		connectionStatus.set('connected');
		reconnectAttempts = 0;
		addLogEvent('connection', 'Connected to server');
	};

	// --- Event: Message Received ---
	socket.onmessage = (event) => {
		try {
			const message = JSON.parse(event.data);
			handleMessage(message);
		} catch (err) {
			console.error('Failed to parse WebSocket message:', err);
		}
	};

	// --- Event: Connection Closed ---
	socket.onclose = (event) => {
		console.log('WebSocket closed:', event.code, event.reason);
		connectionStatus.set('disconnected');

		// Attempt to reconnect if we weren't intentionally closed
		if (currentRoomCode && currentPlayerId) {
			attemptReconnect();
		}
	};

	// --- Event: Connection Error ---
	socket.onerror = (error) => {
		console.error('WebSocket error:', error);
		connectionStatus.set('error');
	};
}

/**
 * Attempts to reconnect after a disconnection.
 * Uses exponential backoff.
 * @private
 */
function attemptReconnect() {
	if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
		console.error('Max reconnection attempts reached');
		addLogEvent('connection', 'Failed to reconnect after multiple attempts');
		return;
	}

	reconnectAttempts++;
	const delay = RECONNECT_DELAY_MS * reconnectAttempts;

	console.log(`Reconnecting in ${delay}ms (attempt ${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS})`);
	connectionStatus.set('reconnecting');
	addLogEvent('connection', `Reconnecting (attempt ${reconnectAttempts})...`);

	reconnectTimer = setTimeout(() => {
		doConnect();
	}, delay);
}

/**
 * Handles incoming WebSocket messages and updates stores.
 * @private
 * @param {{ type: string, data: * }} message
 */
function handleMessage(message) {
	const { type, data } = message;

	switch (type) {
		case 'player_joined':
			handlePlayerJoined(data);
			break;

		case 'player_left':
			handlePlayerLeft(data);
			break;

		case 'phase_changed':
			handlePhaseChanged(data);
			break;

		case 'action_submitted':
			handleActionSubmitted(data);
			break;

		case 'action_deleted':
			handleActionDeleted(data);
			break;

		case 'roles_assigned':
			handleRolesAssigned(data);
			break;

		case 'game_ended':
			handleGameEnded(data);
			break;

		case 'chat':
			handleChat(data);
			break;

		case 'error':
			handleError(data);
			break;

		default:
			console.log('Unknown message type:', type, data);
	}
}

// =============================================================================
// Message Handlers
// =============================================================================

/**
 * @param {{ player_id: string, name?: string }} data
 */
function handlePlayerJoined(data) {
	const playerId = data.player_id;
	const playerName = data.name || 'Unknown';

	room.update((r) => ({
		...r,
		players: [...r.players, { id: playerId, name: playerName, isAlive: true }]
	}));

	addLogEvent('player_joined', `${playerName} joined the room`);
}

/**
 * @param {{ player_id: string }} data
 */
function handlePlayerLeft(data) {
	const playerId = data.player_id;

	room.update((r) => {
		const leavingPlayer = r.players.find((p) => p.id === playerId);
		const name = leavingPlayer?.name || 'Unknown';
		addLogEvent('player_left', `${name} left the room`);

		return {
			...r,
			players: r.players.filter((p) => p.id !== playerId)
		};
	});
}

/**
 * @param {{ phase: 'LOBBY' | 'NIGHT' | 'DAY' }} data
 */
function handlePhaseChanged(data) {
	const newPhase = data.phase;

	room.update((r) => ({ ...r, phase: newPhase }));
	actions.set([]); // Clear actions on phase change

	addLogEvent('phase_changed', `Phase changed to ${newPhase}`);
}

/**
 * @param {{ id: string, playerId: string, type: string, targetId: string | null, phase: string, timestamp: string }} data - Action data
 */
function handleActionSubmitted(data) {
	actions.update((a) => [...a, data]);
	addLogEvent('action_submitted', `Action queued: ${data.type}`);
}

/**
 * @param {{ action_id: string }} data
 */
function handleActionDeleted(data) {
	const actionId = data.action_id;

	actions.update((a) => a.filter((action) => action.id !== actionId));
	addLogEvent('action_deleted', 'Action removed from queue');
}

/**
 * @param {{ role_id: number, role_name: string }} data
 */
function handleRolesAssigned(data) {
	player.update((p) => ({ ...p, roleId: data.role_id, roleName: data.role_name }));
	addLogEvent('roles_assigned', `You are a ${data.role_name}!`);
}

/**
 * @param {{ winner: string, message?: string }} data
 */
function handleGameEnded(data) {
	addLogEvent('game_ended', data.message || `Game over! Winner: ${data.winner}`);
}

/**
 * @param {{ player_id: string, message: string }} data
 */
function handleChat(data) {
	addLogEvent('chat', `${data.player_id}: ${data.message}`);
}

/**
 * @param {{ message: string }} data
 */
function handleError(data) {
	addLogEvent('error', `Error: ${data.message}`);
	console.error('Server error:', data.message);
}
