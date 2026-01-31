<script lang="ts">
  import Avatar from '$lib/components/ui/Avatar.svelte';

  interface Props {
    name: string;
    isHost?: boolean;
    isCurrentPlayer?: boolean;
  }

  let { name, isHost = false, isCurrentPlayer = false }: Props = $props();

  let containerStyles = $derived(isCurrentPlayer
    ? 'bg-primary-500/10 border-2 border-primary-500/30 shadow-md'
    : 'bg-surface-200-800 border-2 border-surface-400-600 shadow-md');
</script>

<div class="flex items-center gap-3 p-2 mb-2 {containerStyles}">
  <Avatar {name} size="md" variant={isCurrentPlayer ? 'primary' : 'surface'} />
  <div class="flex-1 min-w-0">
    <div class="font-medium text-surface-900-50 truncate">
      {name}
    </div>
    <div class="text-xs {isCurrentPlayer ? 'text-primary-600-400' : 'text-surface-600-400'}">
      {isCurrentPlayer ? 'You' : 'Player'} {#if isHost}(Host){/if}
    </div>
  </div>
  {#if isHost}
    <span class="text-lg">👑</span>
  {/if}
</div>
