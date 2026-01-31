// API client for HTTP requests to backend
import type { CreateRoomResponse, JoinRoomResponse, HealthCheckResponse } from './types';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080';

/**
 * Create a new room
 * @param hostName - Name of the host player
 * @returns Promise with room code and host ID
 */
export async function createRoom(hostName: string): Promise<CreateRoomResponse> {
	const response = await fetch(`${API_BASE}/api/rooms`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ hostName })
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({ error: 'Unknown error' }));
		throw new Error(error.error || `Failed to create room: ${response.status}`);
	}

	return response.json();
}

/**
 * Join an existing room
 * @param roomCode - Room code to join
 * @param playerName - Name of the player joining
 * @returns Promise with player ID and phase
 */
export async function joinRoom(roomCode: string, playerName: string): Promise<JoinRoomResponse> {
	const response = await fetch(`${API_BASE}/api/rooms/${roomCode}/join`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ playerName })
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({ error: 'Unknown error' }));
		
		if (response.status === 404) {
			throw new Error('Room not found. Check the code and try again.');
		} else if (response.status === 400) {
			throw new Error('Please check your inputs and try again.');
		}
		
		throw new Error(error.error || `Failed to join room: ${response.status}`);
	}

	return response.json();
}

/**
 * Check database health
 * @returns Promise with health status
 */
export async function healthCheck(): Promise<HealthCheckResponse> {
	const response = await fetch(`${API_BASE}/api/health/db`);
	
	if (!response.ok) {
		throw new Error('Health check failed');
	}
	
	return response.json();
}
