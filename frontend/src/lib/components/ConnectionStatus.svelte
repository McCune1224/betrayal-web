<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';

  interface Props {
    status: 'connected' | 'connecting' | 'disconnected' | 'error' | 'idle';
    error?: string;
    onReconnect: () => void;
  }

  let { status, error, onReconnect }: Props = $props();

  let config = $derived({
    idle: { bg: 'bg-warning-500', text: 'Connecting...', showReconnect: false },
    connecting: { bg: 'bg-warning-500', text: 'Connecting...', showReconnect: false },
    disconnected: { bg: 'bg-error-500', text: 'Disconnected from room', showReconnect: true },
    error: { bg: 'bg-error-500', text: error || 'Unknown error', showReconnect: true },
    connected: { bg: '', text: '', showReconnect: false }
  }[status]);
</script>

{#if status !== 'connected'}
  <div class="fixed top-0 left-0 right-0 z-50 p-3 text-center text-white font-medium {config.bg}">
    <span class={status === 'connecting' ? 'animate-pulse' : ''}>
      {config.text}
    </span>
    {#if config.showReconnect}
      <Button variant="ghost" size="sm" onclick={onReconnect} className="ml-3">
        {status === 'error' ? 'Retry' : 'Reconnect'}
      </Button>
    {/if}
  </div>
{/if}
