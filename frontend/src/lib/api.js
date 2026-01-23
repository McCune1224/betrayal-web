// =============================================================================
// HTTP API Client
// =============================================================================
//
// This module provides functions for communicating with the backend REST API.
// It handles the HTTP requests for creating and joining rooms.
//
// The WebSocket connection is handled separately in ws.js.
//
// =============================================================================

/**
 * Base URL for the backend API.
 * In development: http://localhost:8080
 * In production: Set via VITE_API_BASE environment variable
 */
const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080';

// =============================================================================
// API Response Types (for documentation)
// =============================================================================

/**
 * @typedef {Object} CreateRoomResponse
 * @property {string} code - The room code (e.g., "abc123")
 * @property {string} host - The host's player ID
 */

/**
 * @typedef {Object} JoinRoomResponse
 * @property {string} code - The room code
 * @property {string} player_id - The player's assigned ID
 * @property {string} phase - Current game phase (LOBBY, NIGHT, DAY)
 */

/**
 * @typedef {Object} ApiError
 * @property {string} error - Error message from the server
 */

// =============================================================================
// API Functions
// =============================================================================

/**
 * Creates a new game room.
 * The caller becomes the host of the room.
 *
 * @returns {Promise<CreateRoomResponse>}
 * @throws {Error} If the request fails
 *
 * @example
 * const { code, host } = await createRoom();
 * // code = "abc123", host = "def456..."
 */
export async function createRoom() {
	const response = await fetch(`${API_BASE}/api/rooms`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		}
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to create room');
	}

	return response.json();
}

/**
 * Joins an existing game room.
 *
 * @param {string} code - The room code to join
 * @param {string} name - The player's display name
 * @returns {Promise<JoinRoomResponse>}
 * @throws {Error} If the room doesn't exist or the request fails
 *
 * @example
 * const { player_id, phase } = await joinRoom("abc123", "Alice");
 * // player_id = "xyz789...", phase = "LOBBY"
 */
export async function joinRoom(code, name) {
	const response = await fetch(`${API_BASE}/api/rooms/${code}/join`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ name })
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to join room');
	}

	return response.json();
}

/**
 * Starts the game (host only).
 * Transitions the room from LOBBY to NIGHT phase and assigns roles.
 *
 * @param {string} code - The room code
 * @returns {Promise<void>}
 * @throws {Error} If not the host or the request fails
 */
export async function startGame(code) {
	const response = await fetch(`${API_BASE}/api/rooms/${code}/start`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		}
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to start game');
	}
}

/**
 * Submits a player action for the current phase.
 *
 * @param {string} code - The room code
 * @param {Object} action - The action to submit
 * @param {string} action.type - Action type (e.g., "vote", "kill")
 * @param {string} [action.targetId] - Target player ID (if applicable)
 * @returns {Promise<void>}
 * @throws {Error} If the action is invalid or the request fails
 */
export async function submitAction(code, action) {
	const response = await fetch(`${API_BASE}/api/rooms/${code}/actions`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(action)
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to submit action');
	}
}

/**
 * Deletes an action from the queue (host only).
 *
 * @param {string} code - The room code
 * @param {string} actionId - The action ID to delete
 * @returns {Promise<void>}
 * @throws {Error} If not the host or the request fails
 */
export async function deleteAction(code, actionId) {
	const response = await fetch(`${API_BASE}/api/rooms/${code}/actions/${actionId}`, {
		method: 'DELETE',
		headers: {
			'Content-Type': 'application/json'
		}
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to delete action');
	}
}
