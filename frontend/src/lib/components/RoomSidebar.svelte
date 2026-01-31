<script lang="ts">
  import RoomHeader from './RoomHeader.svelte';
  import PlayerList from './PlayerList.svelte';
  import HostControls from './HostControls.svelte';
  import type { RoomPlayer } from '$lib/types';

  interface Props {
    roomCode: string;
    phase: string;
    players: RoomPlayer[];
    currentPlayerId: string | null;
    currentPlayerName: string;
    isHost: boolean;
    connectionDisabled?: boolean;
    onStartGame?: () => void;
    onAdvancePhase?: () => void;
  }

  let {
    roomCode,
    phase,
    players,
    currentPlayerId,
    currentPlayerName,
    isHost,
    connectionDisabled = false,
    onStartGame,
    onAdvancePhase
  }: Props = $props();
</script>

<div class="w-64 bg-surface-100-900 border-r border-surface-200-800 flex flex-col">
  <RoomHeader {roomCode} {phase} />
  <PlayerList {players} {currentPlayerId} {currentPlayerName} {isHost} />
  <HostControls {phase} {isHost} disabled={connectionDisabled} {onStartGame} {onAdvancePhase} />
</div>
