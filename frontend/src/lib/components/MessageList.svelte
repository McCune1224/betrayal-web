<script lang="ts">
  import ChatMessage from './ChatMessage.svelte';
  import type { Message } from '$lib/types';

  interface Props {
    messages: Message[];
    currentPlayerId: string | null;
  }

  let { messages, currentPlayerId }: Props = $props();

  function formatTime(timestamp: string | undefined): string {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }
</script>

<div class="flex-1 overflow-y-auto p-4 space-y-3">
  {#if messages.length === 0}
    <div class="text-center text-surface-600-400 py-8">
      <p class="text-lg mb-2">Welcome to the room!</p>
      <p class="text-sm">Messages will appear here when players join and chat.</p>
    </div>
  {:else}
    {#each messages as message (message.id)}
      {#if message.type === 'system'}
        <div class="flex justify-center">
          <div class="px-4 py-2 text-sm text-surface-600-400 bg-surface-200-800 border border-surface-400-600 shadow-md">
            {message.text}
            <span class="text-xs opacity-60 ml-2">{formatTime(message.timestamp)}</span>
          </div>
        </div>
      {:else}
        <ChatMessage {message} isOwnMessage={message.senderId === currentPlayerId} />
      {/if}
    {/each}
  {/if}
</div>
