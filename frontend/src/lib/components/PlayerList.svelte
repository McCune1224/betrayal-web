<script lang="ts">
  import PlayerCard from './PlayerCard.svelte';
  import type { RoomPlayer } from '$lib/types';

  interface Props {
    players: RoomPlayer[];
    currentPlayerId: string | null;
    currentPlayerName: string;
    isHost: boolean;
  }

  let { players, currentPlayerId, currentPlayerName, isHost }: Props = $props();

  let otherPlayers = $derived(players.filter(p => p && p.id && p.id !== currentPlayerId));
  let playerCount = $derived(players.length + (currentPlayerId ? 1 : 0));
</script>

<div class="flex-1 overflow-y-auto p-4">
  <div class="text-sm font-medium text-surface-600-400 uppercase tracking-wider mb-3">
    Players ({playerCount})
  </div>

  {#if currentPlayerId}
    <PlayerCard name={currentPlayerName} {isHost} isCurrentPlayer={true} />
  {/if}

  {#each otherPlayers as player}
    <PlayerCard name={player.name} isHost={player.isHost} isCurrentPlayer={false} />
  {/each}
</div>
