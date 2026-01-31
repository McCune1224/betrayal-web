<script lang="ts">
  import Input from '$lib/components/ui/Input.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Alert from '$lib/components/ui/Alert.svelte';

  interface Props {
    username: string;
    roomCode: string;
    error?: string;
    isLoading?: boolean;
    onUsernameChange: (value: string) => void;
    onRoomCodeChange: (value: string) => void;
    onSubmit: () => void;
  }

  let {
    username = $bindable(),
    roomCode = $bindable(),
    error,
    isLoading = false,
    onUsernameChange,
    onRoomCodeChange,
    onSubmit
  }: Props = $props();
</script>

<div class="space-y-5 animate-in fade-in slide-in-from-bottom-2 duration-300">
  <Input
    label="Display Name"
    bind:value={username}
    placeholder="Enter your name..."
    maxlength={20}
    oninput={onUsernameChange}
  />
  <Input
    label="Room Code"
    bind:value={roomCode}
    placeholder="ABC123"
    maxlength={6}
    oninput={onRoomCodeChange}
  />
</div>
<Button onclick={onSubmit} disabled={isLoading} className="w-full">
  {#if isLoading}
    <span class="animate-pulse">Joining...</span>
  {:else}
    Join Room
  {/if}
</Button>
