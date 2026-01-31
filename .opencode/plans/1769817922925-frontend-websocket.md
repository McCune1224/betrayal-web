# Frontend WebSocket Integration Plan

## Overview
Implement real-time chat/communications in the frontend by creating WebSocket integration, state management stores, and the game room page. Focus on getting basic real-time messaging working first.

## Current State
- Landing page exists with create/join forms (UI only, no logic)
- Backend WebSocket and HTTP endpoints are complete
- No stores, API client, or WebSocket client exist yet
- No room page exists

## Goal
Enable real-time communication where multiple players can:
1. Create or join a room
2. See each other join in real-time
3. Send/receive messages via WebSocket
4. See connection status (connected/disconnected)

## Implementation Tasks

### 1. Environment Configuration

**File: `frontend/.env.example`**
Create environment variable template:
```
VITE_API_BASE=http://localhost:8080
VITE_WS_BASE=ws://localhost:8080
```

### 2. Create Stores (`src/lib/stores.js`)

**Purpose:** Centralized state management using Svelte 5 runes

**Stores to create:**
- `player` - { id, name, isHost, joinedAt }
- `room` - { code, phase, players: [], hostId }
- `isHost` - Derived store: player.id === room.hostId
- `connection` - { status: 'connecting' | 'connected' | 'disconnected' | 'error', error: null }
- `messages` - Array of chat/game messages
- `log` - Array of system events

**Features:**
- Use Svelte 5 runes ($state) for reactive stores
- Persist player ID in localStorage for reconnection
- Provide helper functions to update stores

### 3. Create API Client (`src/lib/api.js`)

**Purpose:** HTTP client for REST endpoints

**Functions:**
- `createRoom(hostName)` → POST /api/rooms
  - Returns: { roomCode, hostId }
  - On success: Set player as host, store IDs
  - On error: Throw with user-friendly message

- `joinRoom(roomCode, playerName)` → POST /api/rooms/:code/join
  - Returns: { playerId, phase }
  - On success: Store player ID
  - On error: Handle room not found (404), invalid request (400)

**Features:**
- Base URL from import.meta.env.VITE_API_BASE
- JSON content type headers
- Error handling with specific messages
- Request timeout handling

### 4. Create WebSocket Client (`src/lib/ws.js`)

**Purpose:** WebSocket connection management and message handling

**Functions:**
- `connect(roomCode, playerId, playerName)` 
  - Opens WebSocket connection to /ws?room=CODE&player=ID&name=NAME
  - Sets up event listeners (onopen, onmessage, onclose, onerror)
  - Updates connection store
  - Returns cleanup function

- `disconnect()`
  - Closes WebSocket connection gracefully
  - Updates connection store

- `send(message)`
  - Sends JSON message to server
  - Validates connection state first

- `handleMessage(event)`
  - Parses incoming WebSocket messages
  - Dispatches based on msg.type:
    - `player_joined` → Add player to room.players, add to log
    - `player_left` → Remove player from room.players, add to log
    - `player_rejoined` → Update player status, add to log
    - `phase_changed` → Update room.phase, add to log
    - `error` → Show error notification, log it
    - `pong` → Update connection status (keepalive)

**Features:**
- Automatic reconnection with exponential backoff (optional v1)
- Connection status tracking
- Heartbeat/ping-pong handling
- Error recovery
- Message queue for offline messages (optional v1)

### 5. Create Room Page (`src/routes/room/[code]/+page.svelte`)

**Purpose:** Main game interface where WebSocket communication happens

**On Mount:**
1. Extract room code from URL params
2. Check if player has ID (from store/localStorage)
3. If no player ID, redirect to landing page
4. Connect WebSocket with room code, player ID, player name
5. Set up cleanup on component destroy

**On Destroy:**
- Disconnect WebSocket cleanly
- Clear any intervals/timeouts

**UI Components:**

**Connection Status Banner:**
- Show at top when disconnected or connecting
- Color-coded: green (connected), yellow (connecting), red (disconnected)
- Include reconnect button when disconnected

**Player List Sidebar:**
- Show all players in room
- Indicate host with crown icon
- Indicate current player
- Show connection status per player

**Main Chat Area:**
- Message list showing all events/messages
- System messages (player joined, player left, phase changed)
- Timestamps on messages
- Auto-scroll to latest message

**Message Input (Chat):**
- Input field for typing messages
- Send button
- Handle Enter key to send
- Send "chat_message" type via WebSocket

**Host Controls (visible only to host):**
- "Start Game" button (when in LOBBY phase)
- "Advance Phase" button (when game started)
- Disabled with tooltip when not host

**Debug Panel (optional):**
- Show raw WebSocket messages
- Connection stats
- Useful for development

### 6. Update Landing Page (`src/routes/+page.svelte`)

**Wire up existing UI to actual logic:**

**Create Game flow:**
1. Validate username (non-empty, max 20 chars)
2. Call createRoom(username)
3. On success:
   - Store player as host
   - Navigate to /room/[code]
4. On error: Show error message in UI

**Join Game flow:**
1. Validate username (non-empty, max 20 chars)
2. Validate room code (format: 6 alphanumeric chars uppercase)
3. Call joinRoom(roomCode, username)
4. On success:
   - Store player ID
   - Navigate to /room/[code]
5. On error:
   - Room not found: "Room not found. Check the code and try again."
   - Invalid request: "Please check your inputs and try again."

**Error Display:**
- Show error messages below forms
- Clear error when user starts typing again
- Use Skeleton UI alert components

**Navigation:**
- Use `goto()` from `$app/navigation` to redirect to room page

### 7. Update Layout (`src/routes/+layout.svelte`)

**Add global error handling:**
- Global error boundary for uncaught errors
- Toast notification system for transient errors

### 8. Backend Message Types to Handle

Based on `internal/game/messages.go`, handle these message types:

**Incoming (from server):**
- `player_joined` - New player entered room
- `player_rejoined` - Player reconnected
- `player_left` - Player disconnected
- `phase_changed` - Game phase updated
- `action_submitted` - Player submitted action
- `action_deleted` - Host deleted action
- `roles_assigned` - Game started, roles given
- `game_ended` - Game over, winner announced
- `error` - Server error message
- `pong` - Keepalive response

**Outgoing (to server):**
- `join_room` - Rejoin after reconnect
- `submit_action` - Player action
- `host_command` - Host actions (start game, advance phase, delete action)
- `chat_message` - Chat message
- `ping` - Keepalive

### 9. Real-Time Chat Implementation

**Chat Message Flow:**
1. Player types message and hits Send
2. Client sends `chat_message` WebSocket message with text
3. Server broadcasts to all room members
4. All clients receive and display in chat

**Chat UI:**
- Message bubbles showing sender name
- System messages styled differently (gray, italic)
- Timestamps on hover or always visible
- Auto-scroll to bottom on new messages
- Show "Player is typing..." indicator (optional)

### 10. Connection State Management

**States:**
- `idle` - Not connected yet
- `connecting` - WebSocket opening
- `connected` - Ready for messages
- `disconnected` - Connection closed
- `error` - Connection error

**Transitions:**
- User navigates to room page → `connecting`
- WebSocket onopen → `connected`
- WebSocket onclose → `disconnected`
- WebSocket onerror → `error`
- User clicks reconnect → `connecting`

**UI Feedback:**
- Show banner when not connected
- Disable input when disconnected
- Show reconnection timer/countdown (optional)

## Files to Create/Modify

### New Files:
1. `frontend/.env.example` - Environment template
2. `frontend/src/lib/stores.js` - Svelte stores
3. `frontend/src/lib/api.js` - HTTP client
4. `frontend/src/lib/ws.js` - WebSocket client
5. `frontend/src/routes/room/[code]/+page.svelte` - Room page
6. `frontend/src/lib/components/ConnectionStatus.svelte` - Status banner
7. `frontend/src/lib/components/PlayerList.svelte` - Player list
8. `frontend/src/lib/components/ChatMessage.svelte` - Individual message
9. `frontend/src/lib/components/ChatInput.svelte` - Message input

### Modified Files:
1. `frontend/src/routes/+page.svelte` - Wire up API calls
2. `frontend/svelte.config.js` - Change to adapter-node (for Railway)

## Testing Strategy

### Manual Testing:
1. Start backend server
2. Start frontend dev server
3. Browser 1: Create room, enter username
4. Browser 2: Join room with code, enter different username
5. Verify:
   - Both see each other in player list
   - Chat messages appear on both sides
   - Connection status shows correctly
   - Disconnect one browser, other sees "player left"
   - Reconnect, sees "player rejoined"

### Automated Testing (Future):
- Unit tests for stores
- Unit tests for message parsing
- Integration tests with mock WebSocket

## Success Criteria

- [ ] Can create a room from landing page
- [ ] Can join a room with room code
- [ ] Multiple players see each other in real-time
- [ ] Players see join/leave events
- [ ] Chat messages work between players
- [ ] Connection status is visible
- [ ] Disconnected players can reconnect
- [ ] Host controls visible only to host
- [ ] All errors show user-friendly messages
- [ ] No console errors in browser
- [ ] Works in both light and dark mode

## Rollback Plan

If issues arise:
1. All changes are additive - existing landing page continues to work
2. Can disable WebSocket by not navigating to room page
3. Environment variables control API URLs - can point to mock server

## Next Steps After This

Once real-time chat works:
1. **Phase 4: Game Logic** - Implement actual game mechanics (roles, actions, win conditions)
2. **Phase 5: Deployment** - Docker, Railway deployment
3. **Polish** - Animations, sounds, better error handling
