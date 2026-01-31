<script lang="ts">
  import Input from '$lib/components/ui/Input.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Alert from '$lib/components/ui/Alert.svelte';

  interface Props {
    username: string;
    error?: string;
    isLoading?: boolean;
    onUsernameChange: (value: string) => void;
    onSubmit: () => void;
  }

  let {
    username = $bindable(),
    error,
    isLoading = false,
    onUsernameChange,
    onSubmit
  }: Props = $props();
</script>

<div class="space-y-5 animate-in fade-in slide-in-from-bottom-2 duration-300">
  <div class="p-4 bg-secondary-500/10 border-2 border-secondary-500/30 text-secondary-700-300 text-sm flex items-start gap-3 shadow-md">
    <span class="text-xl">👑</span>
    <p class="leading-relaxed">
      You are creating a new lobby. You will control the game settings and phases.
    </p>
  </div>
  <Input
    label="Host Name"
    bind:value={username}
    placeholder="Enter your name..."
    maxlength={20}
    oninput={onUsernameChange}
  />
</div>
<Button variant="secondary" onclick={onSubmit} disabled={isLoading} className="w-full">
  {#if isLoading}
    <span class="animate-pulse">Creating...</span>
  {:else}
    Create Lobby
  {/if}
</Button>
