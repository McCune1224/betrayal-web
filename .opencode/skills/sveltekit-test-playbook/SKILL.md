---
name: sveltekit-test-playbook
description: Testing guidelines for SvelteKit frontend, including Svelte 5 runes testing patterns and manual testing procedures.
license: MIT
compatibility: opencode
metadata:
  audience: frontend-dev
  pattern: test
---

## What I do
- Provide testing guidance for SvelteKit frontend code
- Manual testing procedures for Svelte 5 runes-based code
- Testing checklists for stores, components, and WebSocket integration
- Debugging tips for common frontend issues

## When to use me
- Before committing frontend changes
- When debugging runes-related issues
- Testing new features or bug fixes
- Verifying store functionality

## Testing Philosophy

This project focuses on **manual testing** rather than automated test frameworks. The priority is:
1. Type safety via `bun run check`
2. Manual browser testing
3. Code review for patterns

---

## Pre-Commit Testing Checklist

### 1. Type Checking (REQUIRED)

Always run before committing:

```bash
cd frontend
bun run check
```

**Must pass with 0 errors.** If there are errors, fix them before committing.

### 2. Svelte 5 Runes Verification

Check for these common issues:

- [ ] Files using `$state`, `$derived`, `$effect` have `.svelte.ts` extension
- [ ] No `$derived` values are exported directly from modules
- [ ] No `$state` primitives are exported (only objects or via getter functions)
- [ ] All imports from stores use `$lib/stores.svelte` (not `$lib/stores`)

### 3. Manual Browser Testing

**Critical paths to test:**

#### Landing Page (`/`)
- [ ] Can switch between "Join Room" and "Create Game" tabs
- [ ] Form validation works (name length, room code format)
- [ ] Error messages display correctly
- [ ] Can create a new room
- [ ] Can join an existing room with valid code
- [ ] Invalid room codes show appropriate errors

#### Room Page (`/room/[code]`)
- [ ] Page loads without errors
- [ ] Room code displays correctly
- [ ] Player list shows current player
- [ ] WebSocket connects (check console for errors)
- [ ] Chat messages can be sent and received
- [ ] Connection status banner appears when disconnected
- [ ] Reconnect button works
- [ ] Host controls appear only for host

### 4. Store Testing

#### Testing Reactive State

When testing stores, verify:

```typescript
// In browser console or component:
import { player, setPlayer, room, addPlayerToRoom } from '$lib/stores.svelte';

// 1. Test state updates
setPlayer('test-id', 'Test Player', true);
console.log(player.name); // Should be 'Test Player'

// 2. Test reactivity
player.name = 'Updated Name';
console.log(player.name); // Should be 'Updated Name'

// 3. Test array operations
addPlayerToRoom({ id: 'p2', name: 'Player 2', isHost: false });
console.log(room.players.length); // Should increase
```

#### Common Store Issues

**Issue**: Changes not reflecting in UI
- **Cause**: Reassigned $state object instead of mutating properties
- **Fix**: Mutate properties: `player.name = 'new'` not `player = {name: 'new'}`

**Issue**: "Cannot export derived state" error
- **Cause**: Exported $derived directly
- **Fix**: Export getter function instead

**Issue**: "rune_outside_svelte" error
- **Cause**: Using runes in .ts file
- **Fix**: Rename to .svelte.ts

---

## WebSocket Testing

### Testing Connection

1. Open browser dev tools → Network → WS (WebSocket)
2. Create or join a room
3. Look for WebSocket connection to `ws://localhost:8080/ws`
4. Verify connection opens without errors

### Testing Message Flow

**Test in two browser windows:**
1. Window A: Create room as Host
2. Window B: Join same room
3. Verify both see each other in player list
4. Send chat message from Window A
5. Verify Window B receives it

**Expected message types in console:**
- `player_joined` - When a player joins
- `player_left` - When a player disconnects
- `chat_message` - Chat messages
- `phase_changed` - When host advances phase

### Common WebSocket Issues

**Issue**: "WebSocket connection failed"
- Check backend is running on port 8080
- Check CORS configuration in backend
- Verify `VITE_WS_BASE` env var is set correctly

**Issue**: Messages not broadcasting
- Check room code matches in both windows
- Verify player IDs are unique
- Check hub is running (see backend logs)

---

## Debugging Tips

### Using Browser DevTools

**Console:**
```javascript
// Access stores directly
const stores = await import('/src/lib/stores.svelte');
console.log(stores.player);
console.log(stores.room);
```

**React Debugger:**
- Use Svelte DevTools extension
- Inspect component state
- Track $state changes

### Common Debug Patterns

**Check if state is reactive:**
```svelte
<script>
  import { player } from '$lib/stores.svelte';
  
  // Add this to see updates
  $effect(() => {
    console.log('Player changed:', player);
  });
</script>
```

**Verify store exports:**
```typescript
// In browser console:
import * as stores from '$lib/stores.svelte';
console.log(Object.keys(stores));
// Should see: player, room, connection, messages, getIsHost, setPlayer, etc.
```

---

## Code Review Checklist

When reviewing frontend PRs, check:

### Svelte 5 Compliance
- [ ] Runes used only in `.svelte.ts` or `.svelte` files
- [ ] No exported $derived values
- [ ] State mutations use property assignment, not reassignment
- [ ] Getter functions used for computed values

### Code Quality
- [ ] Proper TypeScript types on all functions
- [ ] Error handling in async functions
- [ ] No console.log statements in production code
- [ ] Proper cleanup in onDestroy (disconnect WebSocket, clear timers)

### UI/UX
- [ ] Loading states for async actions
- [ ] Error messages are user-friendly
- [ ] Responsive design works on mobile
- [ ] Accessible (ARIA labels, keyboard navigation)

---

## Resources

- **Svelte DevTools**: Browser extension for debugging
- **Project AGENTS.md**: Svelte 5 patterns and conventions
- **sveltekit-ui-patterns skill**: UI implementation patterns
- **Context7 MCP**: Query official Svelte docs for specific errors