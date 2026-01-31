<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import {
    player,
    room,
    connection,
    messages,
    getIsHost,
    loadPlayerFromStorage,
    setRoom,
    addMessage,
  } from "$lib/stores.svelte";
  import { connect, disconnect, sendChatMessage } from "$lib/ws";
  import ConnectionStatus from "$lib/components/ConnectionStatus.svelte";
  import RoomSidebar from "$lib/components/RoomSidebar.svelte";
  import MessageList from "$lib/components/MessageList.svelte";
  import ChatInput from "$lib/components/ChatInput.svelte";

  // Local state
  let chatInput = $state<string>("");

  // Derived state
  let isHost = $derived(getIsHost());
  const roomCode = $derived<string>(page.params.code ?? "");

  onMount(() => {
    if (!roomCode) {
      goto("/");
      return undefined;
    }

    if (!player.id) {
      const loaded = loadPlayerFromStorage();
      if (!loaded) {
        goto("/");
        return undefined;
      }
    }

    setRoom(roomCode);
    const cleanup = connect(roomCode, player.id || "", player.name);

    return () => {
      if (cleanup) cleanup();
    };
  });

  onDestroy(() => {
    disconnect();
  });

  function handleSendMessage(): void {
    if (!chatInput.trim()) return;

    addMessage({
      type: "chat",
      sender: player.name || "Unknown",
      senderId: player.id || "",
      text: chatInput.trim(),
    });

    sendChatMessage(chatInput.trim());
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

  function handleStartGame(): void {
    // TODO: Start game
  }

  function handleAdvancePhase(): void {
    // TODO: Advance phase
  }
</script>

<ConnectionStatus 
  status={connection.status} 
  error={connection.error || undefined} 
  onReconnect={handleReconnect} 
/>

<div
  class="flex-1 flex bg-surface-50-950 h-full overflow-hidden"
  class:pt-12={connection.status !== "connected"}
>
  <RoomSidebar
    {roomCode}
    phase={room.phase}
    players={room.players}
    currentPlayerId={player.id}
    currentPlayerName={player.name}
    {isHost}
    connectionDisabled={connection.status !== "connected"}
    onStartGame={handleStartGame}
    onAdvancePhase={handleAdvancePhase}
  />

  <div class="flex-1 flex flex-col min-w-0">
    <MessageList messages={messages} currentPlayerId={player.id} />
    <ChatInput
      bind:value={chatInput}
      placeholder={connection.status === "connected" ? "Type a message..." : "Connect to chat..."}
      disabled={connection.status !== "connected"}
      onSend={handleSendMessage}
      onKeydown={handleKeydown}
    />
  </div>
</div>
