<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import {
    player,
    room,
    connection,
    getIsHost,
    messages,
    loadPlayerFromStorage,
    setRoom,
    addMessage,
  } from "$lib/stores.svelte";
  import { connect, disconnect, sendChatMessage } from "$lib/ws";

  // Local state
  let chatInput = $state<string>("");
  let messagesContainer = $state<HTMLDivElement | null>(null);

  // Derived state for isHost (reactive)
  let isHost = $derived(getIsHost());

  // Get room code from URL
  const roomCode = $derived<string>(page.params.code ?? "");

  onMount(() => {
    // Guard against undefined room code
    if (!roomCode) {
      goto("/");
      return undefined;
    }

    // Try to load player from storage if not already set
    if (!player.id) {
      const loaded = loadPlayerFromStorage();
      if (!loaded) {
        // No player data, redirect to landing page
        goto("/");
        return undefined;
      }
    }

    // Set room code
    setRoom(roomCode);

    // Connect to WebSocket
    const cleanup = connect(roomCode, player.id || "", player.name);

    return () => {
      if (cleanup) cleanup();
    };
  });

  onDestroy(() => {
    disconnect();
  });

  // Auto-scroll to bottom when new messages arrive
  $effect(() => {
    if (messagesContainer && messages.length > 0) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  });

  function handleSendMessage(): void {
    if (!chatInput.trim()) return;

    // Add message locally first
    addMessage({
      type: "chat",
      sender: player.name || "Unknown",
      senderId: player.id || "",
      text: chatInput.trim(),
    });

    // Send via WebSocket
    sendChatMessage(chatInput.trim());

    // Clear input
    chatInput = "";
  }

  function handleKeydown(event: KeyboardEvent): void {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      handleSendMessage();
    }
  }

  function handleReconnect(): void {
    if (!roomCode) return;
    connect(roomCode, player.id || "", player.name);
  }

  function formatTime(timestamp: string | undefined): string {
    if (!timestamp) return "";
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  }

  function getSenderInitial(sender: string | undefined): string {
    return (sender || "?").charAt(0).toUpperCase();
  }
</script>

<!-- Connection Status Banner -->
{#if connection.status !== "connected"}
  <div
    class="fixed top-0 left-0 right-0 z-50 p-3 text-center text-white font-medium {connection.status ===
    'connecting'
      ? 'bg-warning-500'
      : 'bg-error-500'}"
  >
    {#if connection.status === "connecting"}
      <span class="animate-pulse">Connecting...</span>
    {:else if connection.status === "disconnected"}
      Disconnected from room
      <button
        class="ml-3 px-3 py-1 bg-white/20 hover:bg-white/30 text-sm font-semibold transition-colors shadow"
        onclick={handleReconnect}
      >
        Reconnect
      </button>
    {:else if connection.status === "error"}
      Connection error: {connection.error || "Unknown error"}
      <button
        class="ml-3 px-3 py-1 bg-white/20 hover:bg-white/30 text-sm font-semibold transition-colors shadow"
        onclick={handleReconnect}
      >
        Retry
      </button>
    {/if}
  </div>
{/if}

<!-- Main Layout -->
<div
  class="flex-1 flex bg-surface-50-950 h-full overflow-hidden"
  class:pt-12={connection.status !== "connected"}
>
  <!-- Sidebar - Player List -->
  <div
    class="w-64 bg-surface-100-900 border-r border-surface-200-800 flex flex-col"
  >
    <!-- Room Header -->
    <div class="p-4 border-b border-surface-200-800">
      <div
        class="text-sm font-medium text-surface-600-400 uppercase tracking-wider"
      >
        Room Code
      </div>
      <div
        class="text-2xl font-black font-mono tracking-widest text-surface-900-50"
      >
        {roomCode}
      </div>
      <div class="mt-2 text-sm text-surface-600-400">
        Phase: <span class="font-semibold text-surface-900-50"
          >{room.phase}</span
        >
      </div>
    </div>

    <!-- Player List -->
    <div class="flex-1 overflow-y-auto p-4">
      <div
        class="text-sm font-medium text-surface-600-400 uppercase tracking-wider mb-3"
      >
        Players ({room.players.length + (player.id ? 1 : 0)})
      </div>

      <!-- Current Player -->
      {#if player.id}
        <div
          class="flex items-center gap-3 p-2 bg-primary-500/10 border-2 border-primary-500/30 mb-2 shadow-md"
        >
          <div
            class="w-8 h-8 bg-primary-500 flex items-center justify-center text-white font-bold text-sm shadow"
          >
            {getSenderInitial(player.name)}
          </div>
          <div class="flex-1 min-w-0">
            <div class="font-medium text-surface-900-50 truncate">
              {player.name}
            </div>
            <div class="text-xs text-primary-600-400">
              You {#if isHost}(Host){/if}
            </div>
          </div>
          {#if isHost}
            <span class="text-lg">👑</span>
          {/if}
        </div>
      {/if}

      <!-- Other Players -->
      {#each room.players as otherPlayer}
        {#if otherPlayer && otherPlayer.id && otherPlayer.id !== player.id}
          <div
            class="flex items-center gap-3 p-2 bg-surface-200-800 border-2 border-surface-400-600 mb-2 shadow-md"
          >
            <div
              class="w-8 h-8 bg-surface-400-600 flex items-center justify-center text-white font-bold text-sm shadow"
            >
              {getSenderInitial(otherPlayer.name)}
            </div>
            <div class="flex-1 min-w-0">
              <div class="font-medium text-surface-900-50 truncate">
                {otherPlayer.name}
              </div>
              <div class="text-xs text-surface-600-400">Player</div>
            </div>
            {#if otherPlayer.isHost}
              <span class="text-lg">👑</span>
            {/if}
          </div>
        {/if}
      {/each}
    </div>

    <!-- Host Controls -->
    {#if isHost}
      <div class="p-4 border-t border-surface-200-800 space-y-2">
        <div
          class="text-sm font-medium text-surface-600-400 uppercase tracking-wider mb-2"
        >
          Host Controls
        </div>

        {#if room.phase === "LOBBY"}
          <button
            class="w-full py-2 px-4 font-semibold text-white bg-success-500 hover:bg-success-400 transition-colors border-2 border-success-600 shadow-lg shadow-success-500/30"
            onclick={() => {
              /* TODO: Start game */
            }}
            disabled={connection.status !== "connected"}
          >
            Start Game
          </button>
        {:else}
          <button
            class="w-full py-2 px-4 font-semibold text-white bg-primary-500 hover:bg-primary-400 transition-colors border-2 border-primary-600 shadow-lg shadow-primary-500/30"
            onclick={() => {
              /* TODO: Advance phase */
            }}
            disabled={connection.status !== "connected"}
          >
            Advance Phase
          </button>
        {/if}
      </div>
    {/if}
  </div>

  <!-- Main Chat Area -->
  <div class="flex-1 flex flex-col min-w-0">
    <!-- Messages -->
    <div
      class="flex-1 overflow-y-auto p-4 space-y-3"
      bind:this={messagesContainer}
    >
      {#if messages.length === 0}
        <div class="text-center text-surface-600-400 py-8">
          <p class="text-lg mb-2">Welcome to the room!</p>
          <p class="text-sm">
            Messages will appear here when players join and chat.
          </p>
        </div>
      {:else}
        {#each messages as message (message.id)}
          {#if message.type === "system"}
            <div class="flex justify-center">
              <div
                class="px-4 py-2 text-sm text-surface-600-400 bg-surface-200-800 border border-surface-400-600 shadow-md"
              >
                {message.text}
                <span class="text-xs opacity-60 ml-2"
                  >{formatTime(message.timestamp)}</span
                >
              </div>
            </div>
          {:else}
            <div
              class="flex gap-3 {message.senderId === player.id
                ? 'flex-row-reverse'
                : ''}"
            >
              <div
                class="w-8 h-8 bg-primary-500 flex-shrink-0 flex items-center justify-center text-white font-bold text-sm shadow"
              >
                {getSenderInitial(message.sender)}
              </div>
              <div class="max-w-[70%]">
                <div
                  class="text-xs text-surface-600-400 mb-1 {message.senderId ===
                  player.id
                    ? 'text-right'
                    : ''}"
                >
                  {message.sender || "Unknown"}
                  <span class="opacity-60">{formatTime(message.timestamp)}</span
                  >
                </div>
                <div
                  class="px-4 py-2 border-2 shadow-lg {message.senderId === player.id
                    ? 'bg-primary-500 text-white border-primary-700 border-l-4'
                    : 'bg-surface-200-800 text-surface-900-50 border-surface-400-600 border-r-4'}"
                >
                  {message.text}
                </div>
              </div>
            </div>
          {/if}
        {/each}
      {/if}
    </div>

    <!-- Chat Input -->
    <div class="p-4 border-t border-surface-200-800 bg-surface-100-900">
      <div class="flex gap-2">
        <input
          type="text"
          class="flex-1 px-4 py-3 bg-surface-50-950 text-surface-900-50 border-2 border-surface-300-700 focus:border-primary-500 focus:ring-1 focus:ring-primary-500 outline-none transition-all placeholder:text-surface-400-600 shadow-inner"
          placeholder={connection.status === "connected"
            ? "Type a message..."
            : "Connect to chat..."}
          bind:value={chatInput}
          onkeydown={handleKeydown}
          disabled={connection.status !== "connected"}
        />
        <button
          class="px-6 py-3 font-semibold text-white bg-primary-500 hover:bg-primary-400 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-2 border-primary-600 shadow-lg shadow-primary-500/30"
          onclick={handleSendMessage}
          disabled={!chatInput.trim() || connection.status !== "connected"}
        >
          Send
        </button>
      </div>
    </div>
  </div>
</div>
