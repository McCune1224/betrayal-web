# TypeScript Conversion Plan for Frontend

## Overview
Convert all JavaScript files to TypeScript and add `lang="ts"` to Svelte files that don't already have it. This ensures full type safety and better IDE support throughout the frontend codebase.

## Current State
- 4 JavaScript files need conversion to .ts
- 2 Svelte files need `lang="ts"` added
- JSDoc types already exist and can be converted to TypeScript interfaces
- TypeScript is already configured in the project

## Conversion Tasks

### 1. Convert stores.js to stores.ts

**File:** `frontend/src/lib/stores.js` → `frontend/src/lib/stores.ts`

**Changes:**
- Convert JSDoc @typedef to TypeScript interfaces
- Add proper type annotations to all functions
- Ensure $state, $derived runes are properly typed
- Export types for use in other files

**Types to create:**
```typescript
interface Player {
  id: string | null;
  name: string;
  isHost: boolean;
  joinedAt: string | null;
}

interface RoomPlayer {
  id: string;
  name: string;
  isHost: boolean;
}

interface Room {
  code: string;
  phase: string;
  players: RoomPlayer[];
  hostId: string | null;
}

interface Connection {
  status: 'idle' | 'connecting' | 'connected' | 'disconnected' | 'error';
  error: string | null;
  lastPing: string | null;
}

interface Message {
  id: string;
  type: 'chat' | 'system';
  subtype?: string;
  sender?: string;
  senderId?: string;
  text?: string;
  timestamp: string;
}
```

### 2. Convert ws.js to ws.ts

**File:** `frontend/src/lib/ws.js` → `frontend/src/lib/ws.ts`

**Changes:**
- Import types from stores.ts
- Add types for WebSocket message data
- Type all handler functions properly
- Add proper typing for ws instance and timeout

**Types to add:**
```typescript
interface WSMessage {
  type: string;
  data?: any;
}

interface PlayerJoinedData {
  playerId: string;
  playerName: string;
  isHost: boolean;
}

interface PlayerLeftData {
  playerId: string;
  playerName?: string;
}

interface PhaseChangedData {
  phase: string;
}

interface ActionData {
  playerName: string;
}

interface GameEndedData {
  winner: string;
}

interface ErrorData {
  message: string;
}
```

### 3. Convert api.js to api.ts

**File:** `frontend/src/lib/api.js` → `frontend/src/lib/api.ts`

**Changes:**
- Add return type annotations for all functions
- Define interfaces for API responses
- Add proper error typing

**Types to add:**
```typescript
interface CreateRoomResponse {
  roomCode: string;
  hostId: string;
}

interface JoinRoomResponse {
  playerId: string;
  phase: string;
}

interface HealthCheckResponse {
  status: string;
}
```

### 4. Convert types.js to types.ts

**File:** `frontend/src/lib/types.js` → `frontend/src/lib/types.ts`

**Changes:**
- Convert JSDoc to TypeScript interfaces
- Export all types
- Consider if this file should be merged into stores.ts or kept separate

### 5. Update +page.svelte to use TypeScript

**File:** `frontend/src/routes/+page.svelte`

**Changes:**
- Change `<script>` to `<script lang="ts">`
- Add proper TypeScript types for all variables
- Type the event handlers
- Type API responses

### 6. Update room/[code]/+page.svelte to use TypeScript

**File:** `frontend/src/routes/room/[code]/+page.svelte`

**Changes:**
- Change `<script>` to `<script lang="ts">`
- Import types from stores
- Type all local state variables
- Type event handlers properly
- Ensure proper typing for $page.params

## Import Path Updates

All imports need to be updated from `.js` to the correct path:
- `from './stores.js'` → `from './stores'` (TypeScript resolves without extension)
- OR keep `.js` and let bundler handle it (SvelteKit handles this)

## Files to Modify

### New Files (Renamed from .js to .ts):
1. `frontend/src/lib/stores.ts` (from stores.js)
2. `frontend/src/lib/ws.ts` (from ws.js)
3. `frontend/src/lib/api.ts` (from api.js)
4. `frontend/src/lib/types.ts` (from types.js)

### Modified Files:
1. `frontend/src/routes/+page.svelte` - Add `lang="ts"`
2. `frontend/src/routes/room/[code]/+page.svelte` - Add `lang="ts"`

### Deleted Files:
1. `frontend/src/lib/stores.js`
2. `frontend/src/lib/ws.js`
3. `frontend/src/lib/api.js`
4. `frontend/src/lib/types.js`

## Verification Steps

1. **Type Check:** Run `bun run check` to verify 0 errors
2. **Build:** Run `bun run build` to ensure production build succeeds
3. **Dev Server:** Run `bun dev` to verify dev server starts
4. **Test:** Manually test create/join room flow in browser

## Rollback Plan

If issues arise:
1. Keep original .js files until conversion is verified
2. Can revert by restoring .js files and updating imports
3. All changes are file renames/content updates - no breaking config changes

## Success Criteria

- [ ] All .js files converted to .ts
- [ ] All .svelte files use `<script lang="ts">`
- [ ] `bun run check` shows 0 errors
- [ ] `bun run build` succeeds
- [ ] Dev server starts and runs correctly
- [ ] All functionality works (create room, join room, chat)
