<script>
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { player, room, isHost, log, connectionStatus, resetStores } from '$lib/stores.js';
	import { connectWS, closeWS, sendChat, sendAdvancePhase, sendStartGame } from '$lib/ws.js';

	// Get the room code from the URL
	const roomCode = $derived(page.params.code);

	// Chat input
	let chatInput = $state('');

	// Connect to WebSocket on mount
	onMount(() => {
		const currentPlayer = $player;

		// If we don't have player info, redirect to home
		if (!currentPlayer.id || !roomCode) {
			goto('/');
			return;
		}

		// Connect to WebSocket
		connectWS(roomCode, currentPlayer.id);

		// Cleanup on unmount
		return () => {
			closeWS();
		};
	});

	// Leave room handler
	function handleLeaveRoom() {
		closeWS();
		resetStores();
		goto('/');
	}

	// Send chat message
	function handleSendChat() {
		if (chatInput.trim()) {
			sendChat(chatInput.trim());
			chatInput = '';
		}
	}

	// Handle enter key in chat
	/** @param {KeyboardEvent} event */
	function handleChatKeydown(event) {
		if (event.key === 'Enter') {
			handleSendChat();
		}
	}

	// Start game (host only, from LOBBY)
	function handleStartGame() {
		sendStartGame();
	}

	// Advance phase (host only)
	function handleAdvancePhase() {
		sendAdvancePhase();
	}

	// Get connection status color
	/** @param {string} status */
	function getStatusColor(status) {
		switch (status) {
			case 'connected':
				return 'bg-success-500';
			case 'connecting':
			case 'reconnecting':
				return 'bg-warning-500';
			case 'error':
			case 'disconnected':
				return 'bg-error-500';
			default:
				return 'bg-surface-500';
		}
	}

	// Get phase display info
	/** @param {string | null} phase */
	function getPhaseInfo(phase) {
		switch (phase) {
			case 'LOBBY':
				return { label: 'Lobby', color: 'preset-filled-surface-500' };
			case 'NIGHT':
				return { label: 'Night', color: 'preset-filled-primary-500' };
			case 'DAY':
				return { label: 'Day', color: 'preset-filled-warning-500' };
			default:
				return { label: 'Unknown', color: 'preset-filled-surface-500' };
		}
	}

	// Reactive values
	const currentPlayer = $derived($player);
	const currentRoom = $derived($room);
	const currentIsHost = $derived($isHost);
	const currentLog = $derived($log);
	const currentConnectionStatus = $derived($connectionStatus);
	const phaseInfo = $derived(getPhaseInfo(currentRoom.phase));
</script>

<div class="min-h-screen p-4">
	<div class="max-w-4xl mx-auto space-y-4">
		<!-- Header -->
		<div class="card preset-outlined-surface-200-800 p-4">
			<div class="flex items-center justify-between flex-wrap gap-4">
				<!-- Room Info -->
				<div class="flex items-center gap-4">
					<div>
						<p class="text-sm text-surface-600-400">Room Code</p>
						<p class="text-2xl font-mono font-bold">{roomCode}</p>
					</div>
					<div class="badge {phaseInfo.color}">
						{phaseInfo.label}
					</div>
				</div>

				<!-- Connection Status & Leave -->
				<div class="flex items-center gap-4">
					<div class="flex items-center gap-2">
						<span class="w-3 h-3 rounded-full {getStatusColor(currentConnectionStatus)}"></span>
						<span class="text-sm capitalize">{currentConnectionStatus}</span>
					</div>
					<button
						type="button"
						class="btn btn-sm preset-outlined-error-500"
						onclick={handleLeaveRoom}
					>
						Leave Room
					</button>
				</div>
			</div>
		</div>

		<!-- Main Content Grid -->
		<div class="grid md:grid-cols-3 gap-4">
			<!-- Players List -->
			<div class="card preset-outlined-surface-200-800 p-4 space-y-4">
				<h2 class="text-lg font-semibold">Players ({currentRoom.players.length})</h2>

				{#if currentRoom.players.length === 0}
					<p class="text-surface-600-400 text-sm">Waiting for players to join...</p>
				{:else}
					<ul class="space-y-2">
						{#each currentRoom.players as p (p.id)}
							<li class="flex items-center justify-between p-2 rounded-lg bg-surface-100-900">
								<span class:line-through={!p.isAlive} class:opacity-50={!p.isAlive}>
									{p.name}
								</span>
								{#if p.id === currentRoom.hostId}
									<span class="badge preset-filled-primary-500 text-xs">Host</span>
								{/if}
								{#if p.id === currentPlayer.id}
									<span class="badge preset-filled-secondary-500 text-xs">You</span>
								{/if}
							</li>
						{/each}
					</ul>
				{/if}

				<!-- Your Info -->
				<div class="pt-4 border-t border-surface-300-700">
					<p class="text-sm text-surface-600-400">Playing as</p>
					<p class="font-semibold">{currentPlayer.name || 'Unknown'}</p>
					{#if currentPlayer.roleName}
						<div class="mt-2 p-2 rounded-lg bg-primary-500/20 border border-primary-500">
							<p class="text-xs text-surface-600-400">Your Role</p>
							<p class="font-bold text-primary-500">{currentPlayer.roleName}</p>
						</div>
					{/if}
				</div>
			</div>

			<!-- Event Log & Chat -->
			<div class="md:col-span-2 card preset-outlined-surface-200-800 p-4 space-y-4">
				<h2 class="text-lg font-semibold">Event Log</h2>

				<!-- Log Messages -->
				<div class="h-64 overflow-y-auto space-y-2 p-2 rounded-lg bg-surface-100-900">
					{#if currentLog.length === 0}
						<p class="text-surface-600-400 text-sm text-center py-8">No events yet</p>
					{:else}
						{#each currentLog as event (event.id)}
							<div class="text-sm">
								<span class="text-surface-500">
									{new Date(event.timestamp).toLocaleTimeString()}
								</span>
								<span class="ml-2">{event.message}</span>
							</div>
						{/each}
					{/if}
				</div>

				<!-- Chat Input -->
				<div class="flex gap-2">
					<input
						type="text"
						bind:value={chatInput}
						onkeydown={handleChatKeydown}
						placeholder="Type a message..."
						class="flex-1 px-4 py-2 rounded-lg bg-surface-100-900 border border-surface-300-700 focus:border-primary-500 focus:ring-1 focus:ring-primary-500"
					/>
					<button
						type="button"
						class="btn preset-filled-primary-500"
						onclick={handleSendChat}
					>
						Send
					</button>
				</div>
			</div>
		</div>

		<!-- Host Controls -->
		{#if currentIsHost}
			<div class="card preset-filled-surface-100-900 p-4 space-y-4">
				<h2 class="text-lg font-semibold">Host Controls</h2>

				<div class="flex flex-wrap gap-4">
					{#if currentRoom.phase === 'LOBBY'}
						<button
							type="button"
							class="btn preset-filled-success-500"
							onclick={handleStartGame}
							disabled={currentRoom.players.length < 3}
						>
							Start Game {currentRoom.players.length < 3 ? '(Need 3+ players)' : ''}
						</button>
					{:else}
						<button
							type="button"
							class="btn preset-filled-primary-500"
							onclick={handleAdvancePhase}
						>
							Advance Phase
						</button>
					{/if}
				</div>

				{#if currentRoom.phase === 'LOBBY'}
					<p class="text-sm text-surface-600-400">
						Players: {currentRoom.players.length} (minimum 3 required)
					</p>
				{/if}
			</div>
		{/if}
	</div>
</div>
