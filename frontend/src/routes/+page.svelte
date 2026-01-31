<script lang="ts">
  import { goto } from '$app/navigation';
  import { createRoom, joinRoom } from '$lib/api';
  import { setPlayer, setRoom } from '$lib/stores.svelte';
  import LandingHeader from '$lib/components/LandingHeader.svelte';
  import TabSwitcher from '$lib/components/TabSwitcher.svelte';
  import JoinForm from '$lib/components/JoinForm.svelte';
  import CreateForm from '$lib/components/CreateForm.svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import Alert from '$lib/components/ui/Alert.svelte';
  
  let activeTab = $state<'join' | 'create'>("join");
  let username = $state<string>("");
  let roomCode = $state<string>("");
  let error = $state<string>("");
  let isLoading = $state<boolean>(false);

  function validateUsername(name: string): string | null {
    if (!name || name.trim().length === 0) {
      return "Please enter your name";
    }
    if (name.trim().length > 20) {
      return "Name must be 20 characters or less";
    }
    return null;
  }

  function validateRoomCode(code: string): string | null {
    if (!code || code.trim().length === 0) {
      return "Please enter a room code";
    }
    if (!/^[A-Z0-9]{6}$/i.test(code.trim())) {
      return "Room code must be 6 letters or numbers";
    }
    return null;
  }

  async function joinGame(): Promise<void> {
    error = "";
    
    const usernameError = validateUsername(username);
    if (usernameError) {
      error = usernameError;
      return;
    }

    const codeError = validateRoomCode(roomCode);
    if (codeError) {
      error = codeError;
      return;
    }

    isLoading = true;
    
    try {
      const formattedCode = roomCode.trim().toUpperCase();
      const result = await joinRoom(formattedCode, username.trim());
      
      setPlayer(result.playerId, username.trim(), false);
      setRoom(formattedCode, result.phase);
      
      goto(`/room/${formattedCode}`);
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to join room. Please try again.";
    } finally {
      isLoading = false;
    }
  }

  async function createGame(): Promise<void> {
    error = "";
    
    const usernameError = validateUsername(username);
    if (usernameError) {
      error = usernameError;
      return;
    }

    isLoading = true;
    
    try {
      const result = await createRoom(username.trim());
      
      setPlayer(result.hostId, username.trim(), true);
      setRoom(result.roomCode, 'LOBBY', result.hostId);
      
      goto(`/room/${result.roomCode}`);
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to create room. Please try again.";
    } finally {
      isLoading = false;
    }
  }

  function clearError(): void {
    error = "";
  }

  const tabs = [
    { id: 'join', label: 'Join Room', activeColor: 'primary' },
    { id: 'create', label: 'Create Game', activeColor: 'secondary' }
  ];
</script>

<div class="min-h-full w-full flex flex-col items-center justify-center gap-8 bg-surface-50-950 p-4 transition-colors font-sans">
  <LandingHeader title="Betrayal" tagline="Trust no one." />

  <Card className="w-full max-w-md">
    <TabSwitcher {activeTab} {tabs} onTabChange={(id) => { activeTab = id as 'join' | 'create'; clearError(); }} />
    
    <div class="p-8 space-y-6">
      {#if error}
        <Alert type="error" message={error} />
      {/if}
      
      {#if activeTab === "join"}
        <JoinForm
          bind:username
          bind:roomCode
          {isLoading}
          onUsernameChange={clearError}
          onRoomCodeChange={clearError}
          onSubmit={joinGame}
        />
      {:else}
        <CreateForm
          bind:username
          {isLoading}
          onUsernameChange={clearError}
          onSubmit={createGame}
        />
      {/if}
    </div>
  </Card>
</div>
