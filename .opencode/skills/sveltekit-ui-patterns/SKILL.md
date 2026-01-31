---
name: sveltekit-ui-patterns
description: Guidance for SvelteKit component/page structure, error handling, and component composition using Skeleton UI. Includes comprehensive Svelte 5 runes documentation.
license: MIT
compatibility: opencode
metadata:
  audience: frontend-dev
  pattern: ui-best-practice
---

## What I do
- Provide standard recipes for SvelteKit page/component/store structure, and integration with Skeleton UI for layouts and modals
- Patterns for error/documentation blocks, container/layout usage, and accessibility (ARIA/role)
- Examples for validation, user feedback, and Svelte 5 runes usage
- Comprehensive guide for Svelte 5 runes patterns and common pitfalls

## When to use me
Reference this skill when:
- Creating or reviewing SvelteKit page/components
- Working with Svelte 5 runes ($state, $derived, $effect)
- Sharing state between modules
- Debugging runes-related errors
- Keeping UI, error handling, and structure consistent

## Example Usage
Used in `frontend/src/routes/+page.svelte`, new `src/routes/room/[code]/+page.svelte`, when adopting Skeleton UI, or when implementing Svelte 5 runes patterns.

---

## Svelte 5 Runes - Quick Reference

### File Extensions Matter

**CRITICAL RULE**: Files using runes MUST use `.svelte.ts` or `.svelte.js` extension.

| Extension | Can Use Runes? | Example |
|-----------|---------------|---------|
| `.svelte.ts` | ✅ Yes | `stores.svelte.ts` |
| `.svelte.js` | ✅ Yes | `utils.svelte.js` |
| `.svelte` | ✅ Yes | Components |
| `.ts` | ❌ No | Regular TypeScript |
| `.js` | ❌ No | Regular JavaScript |

**Common Error**: Using `$state` in a `.ts` file throws `rune_outside_svelte` error.

### Export Rules - The Complete Guide

#### ❌ NEVER Export These

**1. $derived values directly:**
```typescript
// ❌ WRONG - Will throw "derived_invalid_export" error
export const isHost = $derived(Boolean(player.id && room.hostId && player.id === room.hostId));

// ✅ CORRECT - Export getter function
export function getIsHost(): boolean {
  return Boolean(player.id && room.hostId && player.id === room.hostId);
}
```

**2. $state primitives that get reassigned:**
```typescript
// ❌ WRONG - Breaks reactivity in importing modules
export let count = $state(0);

// ✅ CORRECT - Keep private, export accessors
let count = $state(0);
export function getCount() { return count; }
export function increment() { count += 1; }
export function decrement() { count -= 1; }
```

#### ✅ SAFE to Export

**1. $state objects (mutate properties, don't reassign):**
```typescript
// ✅ CORRECT - Export object, mutate properties
export const player = $state<Player>({
  id: null,
  name: '',
  isHost: false,
  joinedAt: null
});

// ✅ This works - mutating property
player.name = 'New Name';

// ❌ AVOID - Reassigning breaks reactivity
player = { name: 'New Name', id: '123', isHost: false, joinedAt: null };
```

**2. Getter functions for derived values:**
```typescript
// ✅ CORRECT - Export getter
export function getIsHost(): boolean {
  return Boolean(player.id && room.hostId && player.id === room.hostId);
}

// ✅ CORRECT - Export multiple getters
export function getPlayerCount(): number {
  return room.players.length + 1;
}

export function getCurrentPhase(): string {
  return room.phase;
}
```

### Pattern: Complete Store Module

```typescript
// File: stores.svelte.ts
import type { Player, Room, Connection, Message } from './types';

// ✅ Export $state objects (safe)
export const player = $state<Player>({
  id: null,
  name: '',
  isHost: false,
  joinedAt: null
});

export const room = $state<Room>({
  code: '',
  phase: 'LOBBY',
  players: [],
  hostId: null
});

export const connection = $state<Connection>({
  status: 'idle',
  error: null,
  lastPing: null
});

export const messages = $state<Message[]>([]);

// ✅ Export getter functions for derived state
export function getIsHost(): boolean {
  return Boolean(player.id && room.hostId && player.id === room.hostId);
}

export function getPlayerCount(): number {
  return room.players.length + (player.id ? 1 : 0);
}

// ✅ Export action functions
export function setPlayer(id: string, name: string, isHost: boolean = false): void {
  player.id = id;
  player.name = name;
  player.isHost = isHost;
  player.joinedAt = new Date().toISOString();
}

export function addPlayerToRoom(playerData: RoomPlayer): void {
  const existingIndex = room.players.findIndex(p => p.id === playerData.id);
  if (existingIndex >= 0) {
    room.players[existingIndex] = { ...room.players[existingIndex], ...playerData };
  } else {
    room.players.push(playerData);
  }
}

export function setConnectionStatus(status: Connection['status']): void {
  connection.status = status;
}
```

### Pattern: Usage in Components

```svelte
<!-- File: +page.svelte -->
<script lang="ts">
  import { 
    player, 
    room, 
    connection,
    getIsHost,
    setPlayer,
    addPlayerToRoom
  } from '$lib/stores.svelte';
  
  // ✅ Create local derived state in component
  let isHost = $derived(getIsHost());
  let playerCount = $derived(room.players.length + 1);
  
  function handleJoin() {
    // ✅ Call action functions
    setPlayer('player-123', 'Alice', false);
  }
</script>

<div>
  <h1>Welcome {player.name}</h1>
  
  {#if isHost}
    <span class="badge">Host</span>
  {/if}
  
  <p>Room: {room.code}</p>
  <p>Players: {playerCount}</p>
  <p>Status: {connection.status}</p>
</div>
```

### Common Errors and Solutions

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `rune_outside_svelte` | Using `$state` in `.ts` file | Rename file to `.svelte.ts` |
| `derived_invalid_export` | Exporting `$derived` directly | Export getter function instead |
| `state_export_reassigned` | Exporting reassigned `$state` | Use object state or getter |

### Why These Restrictions Exist

The Svelte compiler transforms `$state` references by wrapping them in `$.get()` and `$.set()` calls. The compiler operates on **one file at a time**. If another file imports a directly reassigned `$state` variable, the compiler doesn't know to apply these transformations in the importing module, leading to:

- Stale values
- Broken reactivity  
- Runtime errors

**Rule of thumb**: Export objects you mutate OR export functions, never export reassigned primitives or derived values.

### Pre-Commit Checklist

Before committing frontend changes:

1. ✅ Run `bun run check` - Must pass with 0 errors
2. ✅ Verify all files using runes have `.svelte.ts` extension
3. ✅ Verify no `$derived` values are exported directly
4. ✅ Verify no `$state` primitives are exported (use objects or getters)
5. ✅ Test the application manually in browser

---

## Component Patterns

### Error Handling Pattern

```svelte
<script lang="ts">
  let error = $state<string>("");
  let isLoading = $state<boolean>(false);
  
  async function handleAction() {
    error = "";
    isLoading = true;
    
    try {
      await someAsyncAction();
    } catch (err) {
      error = err instanceof Error ? err.message : "An error occurred";
    } finally {
      isLoading = false;
    }
  }
</script>

{#if error}
  <div class="error-banner">
    {error}
  </div>
{/if}

<button onclick={handleAction} disabled={isLoading}>
  {#if isLoading}
    <span class="animate-pulse">Loading...</span>
  {:else}
    Submit
  {/if}
</button>
```

### Form Validation Pattern

```svelte
<script lang="ts">
  let username = $state<string>("");
  let validationError = $state<string | null>(null);
  
  function validateUsername(name: string): string | null {
    if (!name || name.trim().length === 0) {
      return "Please enter your name";
    }
    if (name.trim().length > 20) {
      return "Name must be 20 characters or less";
    }
    return null;
  }
  
  function handleSubmit() {
    validationError = validateUsername(username);
    if (!validationError) {
      // Proceed with submission
    }
  }
</script>

<input 
  type="text" 
  bind:value={username}
  oninput={() => validationError = null}
  class:error={validationError}
/>

{#if validationError}
  <span class="error-text">{validationError}</span>
{/if}
```

---

## Skeleton UI Integration

### Theme Classes

This project uses Skeleton UI's theme classes (e.g., `bg-surface-50-950`, `text-surface-900-50`). These automatically adapt to light/dark mode.

Common patterns:
- `bg-surface-50-950` - Background (light: 50, dark: 950)
- `text-surface-900-50` - Text (light: 900, dark: 50)
- `border-surface-200-800` - Borders
- `bg-primary-500` - Primary brand color
- `bg-error-500` - Error states

### Layout Structure

```svelte
<div class="min-h-full w-full flex flex-col items-center justify-center 
            bg-surface-50-950 p-4">
  <div class="w-full max-w-md rounded-2xl bg-surface-100-900 
              shadow-2xl border border-surface-200-800">
    <!-- Content -->
  </div>
</div>
```

---

## Resources

- **Svelte 5 Runes Docs**: https://svelte.dev/docs/svelte/what-are-runes
- **Svelte 5 State Docs**: https://svelte.dev/docs/svelte/$state
- **SvelteKit Docs**: https://svelte.dev/docs/kit
- **Skeleton UI**: https://www.skeleton.dev/
- **Project AGENTS.md**: Reference for project-specific conventions
- **Context7 MCP**: Use for up-to-date Svelte documentation queries