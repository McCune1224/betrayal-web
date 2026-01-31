<script lang="ts">
  interface Props {
    roomCode: string;
    phase: string;
  }

  let { roomCode, phase }: Props = $props();
</script>

<div class="p-4 border-b border-surface-200-800">
  <div class="text-sm font-medium text-surface-600-400 uppercase tracking-wider">
    Room Code
  </div>
  <div class="text-2xl font-black font-mono tracking-widest text-surface-900-50">
    {roomCode}
  </div>
  <div class="mt-2 text-sm text-surface-600-400">
    Phase: <span class="font-semibold text-surface-900-50">{phase}</span>
  </div>
</div>
