// =============================================================================
// Svelte Stores for Game State
// =============================================================================
//
// Stores are Svelte's reactive state containers. They allow components to
// subscribe to changes and automatically re-render when data updates.
//
// We use writable stores for state that can change, and derived stores for
// computed values that depend on other stores.
//
// =============================================================================

import { writable, derived } from 'svelte/store';

// =============================================================================
// Player Store
// =============================================================================

/**
 * Current player's information.
 * This is set when joining a room and persists for the session.
 *
 * @type {import('svelte/store').Writable<{
 *   id: string | null,
 *   name: string | null,
 *   roleId: number | null,
 *   roleName: string | null,
 *   isAlive: boolean
 * }>}
 */
export const player = writable({
	id: null,
	name: null,
	roleId: null,
	roleName: null,
	isAlive: true
});

// =============================================================================
// Room Store
// =============================================================================

/**
 * Current room state.
 * Updated via WebSocket messages from the server.
 *
 * @type {import('svelte/store').Writable<{
 *   code: string | null,
 *   hostId: string | null,
 *   phase: 'LOBBY' | 'NIGHT' | 'DAY' | null,
 *   players: Array<{id: string, name: string, isAlive: boolean}>
 * }>}
 */
export const room = writable({
	code: null,
	hostId: null,
	phase: null,
	players: []
});

// =============================================================================
// Derived Stores
// =============================================================================

/**
 * Whether the current player is the room host.
 * Hosts have special privileges like advancing phases.
 */
export const isHost = derived([player, room], ([$player, $room]) => {
	return $player.id !== null && $player.id === $room.hostId;
});

// =============================================================================
// Actions Store
// =============================================================================

/**
 * Queue of pending actions for the current phase.
 * Actions are submitted by players and resolved when the phase advances.
 *
 * @type {import('svelte/store').Writable<Array<{
 *   id: string,
 *   playerId: string,
 *   type: string,
 *   targetId: string | null,
 *   phase: string,
 *   timestamp: string
 * }>>}
 */
export const actions = writable([]);

// =============================================================================
// Event Log Store
// =============================================================================

/**
 * Log of game events for display in the UI.
 * New events are added to the end of the array.
 *
 * @type {import('svelte/store').Writable<Array<{
 *   id: string,
 *   type: string,
 *   message: string,
 *   timestamp: string
 * }>>}
 */
export const log = writable([]);

// =============================================================================
// Connection Status Store
// =============================================================================

/**
 * WebSocket connection status for UI feedback.
 */
export const connectionStatus = writable('disconnected');

// =============================================================================
// Helper Functions
// =============================================================================

/**
 * Adds an event to the log.
 * @param {string} type - Event type (e.g., 'player_joined', 'phase_changed')
 * @param {string} message - Human-readable message
 */
export function addLogEvent(type, message) {
	log.update((events) => [
		...events,
		{
			id: crypto.randomUUID(),
			type,
			message,
			timestamp: new Date().toISOString()
		}
	]);
}

/**
 * Resets all stores to initial state.
 * Call this when leaving a room or on error.
 */
export function resetStores() {
	player.set({ id: null, name: null, roleId: null, roleName: null, isAlive: true });
	room.set({ code: null, hostId: null, phase: null, players: [] });
	actions.set([]);
	log.set([]);
	connectionStatus.set('disconnected');
}
