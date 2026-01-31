<script lang="ts">
  import Avatar from '$lib/components/ui/Avatar.svelte';
  import type { Message } from '$lib/types';

  interface Props {
    message: Message;
    isOwnMessage: boolean;
  }

  let { message, isOwnMessage }: Props = $props();

  function formatTime(timestamp: string | undefined): string {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  let messageStyles = $derived(isOwnMessage
    ? 'bg-primary-500 text-white border-primary-700 border-l-4'
    : 'bg-surface-200-800 text-surface-900-50 border-surface-400-600 border-r-4');
</script>

<div class="flex gap-3 {isOwnMessage ? 'flex-row-reverse' : ''}">
  <Avatar name={message.sender || 'Unknown'} size="md" variant="primary" />
  <div class="max-w-[70%]">
    <div class="text-xs text-surface-600-400 mb-1 {isOwnMessage ? 'text-right' : ''}">
      {message.sender || 'Unknown'}
      <span class="opacity-60">{formatTime(message.timestamp)}</span>
    </div>
    <div class="px-4 py-2 border-2 shadow-lg {messageStyles}">
      {message.text}
    </div>
  </div>
</div>
