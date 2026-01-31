<script lang="ts">
  import { goto } from '$app/navigation';
  import { createRoom, joinRoom } from '$lib/api';
  import { setPlayer, setRoom } from '$lib/stores.svelte';
  
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
    // Room codes are 6 alphanumeric characters
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
      
      // Set player and room in stores
      setPlayer(result.playerId, username.trim(), false);
      setRoom(formattedCode, result.phase);
      
      // Navigate to room page
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
      
      // Set player as host
      setPlayer(result.hostId, username.trim(), true);
      setRoom(result.roomCode, 'LOBBY', result.hostId);
      
      // Navigate to room page
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
</script>

<!-- Main Container -->
<div
  class="min-h-full w-full flex flex-col items-center justify-center gap-8 bg-surface-50-950 p-4 transition-colors font-sans"
>
  <!-- Header -->
  <div class="text-center space-y-2">
    <h1 class="text-5xl font-black tracking-tight text-surface-900-50">
      Betrayal
    </h1>
    <p
      class="text-xl font-medium text-surface-600-400 tracking-wide opacity-80"
    >
      Trust no one.
    </p>
  </div>

  <!-- Card -->
  <div
    class="w-full max-w-md bg-surface-100-900 shadow-2xl border-2 border-surface-300-700 overflow-hidden backdrop-blur-sm"
  >
    <!-- Tabs -->
    <div class="flex border-b border-surface-200-800 relative">
      <button
        class="flex-1 py-4 text-center font-bold tracking-wide text-sm uppercase transition-all duration-200 {activeTab ===
        'join'
          ? 'bg-surface-100-900 text-primary-600-400 shadow-[inset_0_-2px_0_0_rgba(0,0,0,0)] border-b-2 border-primary-500'
          : 'bg-surface-200-800 text-surface-500 hover:text-surface-700-300 hover:bg-surface-300-700'}"
        onclick={() => { activeTab = "join"; clearError(); }}
      >
        Join Room
      </button>
      <button
        class="flex-1 py-4 text-center font-bold tracking-wide text-sm uppercase transition-all duration-200 {activeTab ===
        'create'
          ? 'bg-surface-100-900 text-secondary-600-400 border-b-2 border-secondary-500'
          : 'bg-surface-200-800 text-surface-500 hover:text-surface-700-300 hover:bg-surface-300-700'}"
        onclick={() => { activeTab = "create"; clearError(); }}
      >
        Create Game
      </button>
    </div>

    <!-- Content Area -->
    <div class="p-8 space-y-6">
      <!-- Error Message -->
      {#if error}
        <div class="p-4 bg-error-500/10 border-2 border-error-500/30 text-error-700-300 text-sm shadow-md">
          {error}
        </div>
      {/if}
      
      {#if activeTab === "join"}
        <div
          class="space-y-5 animate-in fade-in slide-in-from-bottom-2 duration-300"
        >
          <label class="block space-y-2">
            <span
              class="text-sm font-bold uppercase tracking-wider text-surface-600-400"
              >Display Name</span
            >
            <input
              class="w-full px-4 py-3 bg-surface-50-950 text-surface-900-50 border-2 border-surface-300-700 focus:border-primary-500 focus:ring-1 focus:ring-primary-500 outline-none transition-all placeholder:text-surface-400-600 shadow-inner"
              type="text"
              placeholder="Enter your name..."
              bind:value={username}
              oninput={clearError}
              maxlength="20"
            />
          </label>
          <label class="block space-y-2">
            <span
              class="text-sm font-bold uppercase tracking-wider text-surface-600-400"
              >Room Code</span
            >
            <input
              class="w-full px-4 py-3 bg-surface-50-950 text-surface-900-50 border-2 border-surface-300-700 focus:border-primary-500 focus:ring-1 focus:ring-primary-500 outline-none transition-all uppercase tracking-[0.2em] font-mono text-center placeholder:text-surface-400-600 shadow-inner"
              type="text"
              placeholder="ABC123"
              maxlength="6"
              bind:value={roomCode}
              oninput={clearError}
            />
          </label>
        </div>
        <button
          class="w-full py-3.5 font-bold text-white bg-gradient-to-br from-primary-500 to-primary-600 hover:from-primary-400 hover:to-primary-500 shadow-lg shadow-primary-500/30 transition-all hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed border-2 border-primary-700"
          onclick={joinGame}
          disabled={isLoading}
        >
          {#if isLoading}
            <span class="animate-pulse">Joining...</span>
          {:else}
            Join Room
          {/if}
        </button>
      {:else}
        <div
          class="space-y-5 animate-in fade-in slide-in-from-bottom-2 duration-300"
        >
          <div
            class="p-4 bg-secondary-500/10 border-2 border-secondary-500/30 text-secondary-700-300 text-sm flex items-start gap-3 shadow-md"
          >
            <span class="text-xl">👑</span>
            <p class="leading-relaxed">
              You are creating a new lobby. You will control the game settings
              and phases.
            </p>
          </div>
          <label class="block space-y-2">
            <span
              class="text-sm font-bold uppercase tracking-wider text-surface-600-400"
              >Host Name</span
            >
            <input
              class="w-full px-4 py-3 bg-surface-50-950 text-surface-900-50 border-2 border-surface-300-700 focus:border-secondary-500 focus:ring-1 focus:ring-secondary-500 outline-none transition-all placeholder:text-surface-400-600 shadow-inner"
              type="text"
              placeholder="Enter your name..."
              bind:value={username}
              oninput={clearError}
              maxlength="20"
            />
          </label>
        </div>
        <button
          class="w-full py-3.5 font-bold text-white bg-gradient-to-br from-secondary-500 to-secondary-600 hover:from-secondary-400 hover:to-secondary-500 shadow-lg shadow-secondary-500/30 transition-all hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed border-2 border-secondary-700"
          onclick={createGame}
          disabled={isLoading}
        >
          {#if isLoading}
            <span class="animate-pulse">Creating...</span>
          {:else}
            Create Lobby
          {/if}
        </button>
      {/if}
    </div>
  </div>

  <!-- Footer -->
  <!-- <div class="text-surface-400-600 text-xs"> -->
  <!--   <p>&copy; {new Date().getFullYear()}</p> -->
  <!-- </div> -->
</div>
