<script lang="ts">
  let activeTab = $state("join");
  let username = $state("");
  let roomCode = $state("");

  function joinGame() {
    console.log("Join", { username, roomCode });
  }

  function createGame() {
    console.log("Create", { username });
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
    class="w-full max-w-md rounded-2xl bg-surface-100-900 shadow-2xl border border-surface-200-800 overflow-hidden backdrop-blur-sm"
  >
    <!-- Tabs -->
    <div class="flex border-b border-surface-200-800 relative">
      <button
        class="flex-1 py-4 text-center font-bold tracking-wide text-sm uppercase transition-all duration-200 {activeTab ===
        'join'
          ? 'bg-surface-100-900 text-primary-600-400 shadow-[inset_0_-2px_0_0_rgba(0,0,0,0)] border-b-2 border-primary-500'
          : 'bg-surface-200-800 text-surface-500 hover:text-surface-700-300 hover:bg-surface-300-700'}"
        onclick={() => (activeTab = "join")}
      >
        Join Room
      </button>
      <button
        class="flex-1 py-4 text-center font-bold tracking-wide text-sm uppercase transition-all duration-200 {activeTab ===
        'create'
          ? 'bg-surface-100-900 text-secondary-600-400 border-b-2 border-secondary-500'
          : 'bg-surface-200-800 text-surface-500 hover:text-surface-700-300 hover:bg-surface-300-700'}"
        onclick={() => (activeTab = "create")}
      >
        Create Game
      </button>
    </div>

    <!-- Content Area -->
    <div class="p-8 space-y-6">
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
              class="w-full px-4 py-3 rounded-lg bg-surface-50-950 text-surface-900-50 border border-surface-300-700 focus:border-primary-500 focus:ring-1 focus:ring-primary-500 outline-none transition-all placeholder:text-surface-400-600"
              type="text"
              placeholder="Enter your name..."
              bind:value={username}
            />
          </label>
          <label class="block space-y-2">
            <span
              class="text-sm font-bold uppercase tracking-wider text-surface-600-400"
              >Room Code</span
            >
            <input
              class="w-full px-4 py-3 rounded-lg bg-surface-50-950 text-surface-900-50 border border-surface-300-700 focus:border-primary-500 focus:ring-1 focus:ring-primary-500 outline-none transition-all uppercase tracking-[0.2em] font-mono text-center placeholder:text-surface-400-600"
              type="text"
              placeholder="ABCD"
              maxlength="4"
              bind:value={roomCode}
            />
          </label>
        </div>
        <button
          class="w-full py-3.5 rounded-lg font-bold text-white bg-gradient-to-br from-primary-500 to-primary-600 hover:from-primary-400 hover:to-primary-500 shadow-lg shadow-primary-500/30 transition-all hover:scale-[1.02] active:scale-[0.98]"
          onclick={joinGame}
        >
          Join Room
        </button>
      {:else}
        <div
          class="space-y-5 animate-in fade-in slide-in-from-bottom-2 duration-300"
        >
          <div
            class="p-4 rounded-lg bg-secondary-500/10 border border-secondary-500/20 text-secondary-700-300 text-sm flex items-start gap-3"
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
              class="w-full px-4 py-3 rounded-lg bg-surface-50-950 text-surface-900-50 border border-surface-300-700 focus:border-secondary-500 focus:ring-1 focus:ring-secondary-500 outline-none transition-all placeholder:text-surface-400-600"
              type="text"
              placeholder="Enter your name..."
              bind:value={username}
            />
          </label>
        </div>
        <button
          class="w-full py-3.5 rounded-lg font-bold text-white bg-gradient-to-br from-secondary-500 to-secondary-600 hover:from-secondary-400 hover:to-secondary-500 shadow-lg shadow-secondary-500/30 transition-all hover:scale-[1.02] active:scale-[0.98]"
          onclick={createGame}
        >
          Create Lobby
        </button>
      {/if}
    </div>
  </div>

  <!-- Footer -->
  <!-- <div class="text-surface-400-600 text-xs"> -->
  <!--   <p>&copy; {new Date().getFullYear()}</p> -->
  <!-- </div> -->
</div>
