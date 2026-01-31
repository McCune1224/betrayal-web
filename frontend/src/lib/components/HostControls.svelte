<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';

  interface Props {
    phase: string;
    isHost: boolean;
    disabled?: boolean;
    onStartGame?: () => void;
    onAdvancePhase?: () => void;
  }

  let { phase, isHost, disabled = false, onStartGame, onAdvancePhase }: Props = $props();
</script>

{#if isHost}
  <div class="p-4 border-t border-surface-200-800 space-y-2">
    <div class="text-sm font-medium text-surface-600-400 uppercase tracking-wider mb-2">
      Host Controls
    </div>

    {#if phase === 'LOBBY'}
      <Button variant="success" onclick={onStartGame} {disabled} className="w-full">
        Start Game
      </Button>
    {:else}
      <Button variant="primary" onclick={onAdvancePhase} {disabled} className="w-full">
        Advance Phase
      </Button>
    {/if}
  </div>
{/if}
