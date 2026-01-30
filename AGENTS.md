# AI Agent Guidelines for Social Deduction Game PoC

---

## Agent Skills Usage Policy (AI/human guidance)

This project uses [OpenCode Agent Skills](https://opencode.ai/docs/skills) for reusable, role-specific, step-by-step knowledge. These are stored in `.opencode/skills/` as SKILL.md files. **Any agent (AI or developer) MUST check and follow relevant skills before implementing, refactoring, or reviewing code, or onboarding others.**

### What are skills?
- Each is a self-contained markdown file describing a workflow, best practice, or checklist for part of the stack (e.g., WebSocket Go handler, SvelteKit UI, DB migration, CI/CD, error handling, game logic flow, etc).
- Skills explain **how to do something correctly in this repo**, referencing actual component structure, dependencies, and conventions.

### When/how to use them
- **Before starting work, always check which skills may be relevant.**
- You can read skills manually or use the agent’s `skill` loading feature to auto-apply them before planning or coding.
- When requesting agent or human help, always mention which skills were referenced (or attach their titles to your PR/issue).
- If you’re unsure which skills to use, list the available skills in `.opencode/skills/` and ask a maintainer or agent.

### Permissions
- All skills are MIT-licensed; internal convention is that skill usage is **always allowed** for all contributors and AI participants.
- The root `opencode.json` config (see below) must specify (or default to) `{ "skill": { "*": "allow" } }` to enable access for all.

---

> Check `CURRENT-STATE.md` for what actually exists before assuming code is present.

Use this file to direct AI agents (Claude, OpenAI, etc.) for efficient project development. Copy-paste sections into your editor or AI tool as needed.

---

## Quick Links

- **Current State:** `CURRENT-STATE.md` - What actually exists now (ground truth)
- **MVP Goals:** `MVP-GOALS.md` - Prioritized checklist with testing patterns
- **Project Brief:** `project-brief.md` - Requirements and specifications
- **Backend:** `backend/`
- **Frontend:** `frontend/`

---

## Project Context

**Name:** Betrayal
**Stack:** Go (Echo) + SvelteKit (Bun) + Postgres + Railway  
**Status:** Early scaffold complete; WebSocket handlers + frontend integration TODO  

**Tech Details:**
- Backend: Go 1.25, Echo framework, gorilla/websocket, pgx + sqlc
- Frontend: SvelteKit, Bun runtime, Skeleton UI (optional), native WebSocket API
- Database: Postgres 14+, golang-migrate, sqlc code generation
- Deployment: Railway (separate services for backend, frontend, Postgres)

---

## Architecture Quick Reference

### Backend Structure
```
backend/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── db.go                        # Database connection (pgx)
│   ├── models/types.go              # Game data structures
│   ├── game/
│   │   ├── room.go                  # Room/player management
│   │   ├── hub.go                   # WebSocket hub & broadcasting
│   │   └── errors.go                # Game error types
│   └── handlers/
│       ├── rooms.go                 # HTTP: POST /api/rooms, /api/rooms/:code/join
│       └── ws.go                    # WebSocket: GET /ws (TODO: implement)
├── internal/db/
│   ├── migrations/                  # golang-migrate SQL files
│   └── sql/                         # sqlc query definitions
├── sqlc.yaml                        # sqlc configuration
├── Dockerfile                       # Railway build
└── .env.example / .env              # Database & port config
```

### Frontend Structure
```
frontend/
├── src/
│   ├── routes/
│   │   ├── +page.svelte             # Landing (create/join room)
│   │   └── room/[code]/+page.svelte # In-game page
│   ├── lib/
│   │   ├── stores.js                # Svelte stores for game state
│   │   ├── api.js                   # HTTP client (createRoom, joinRoom)
│   │   └── ws.js                    # WebSocket client (connectWS, sendMessage)
│   └── app.html
├── svelte.config.js                 # adapter-node for Railway
├── Dockerfile                       # Railway build
└── .env.example / .env              # Backend API/WS URLs
```

### Game State Model
```
Room {
  code: string                    # Room identifier
  hostID: string                  # Out-of-game host (can't win)
  phase: "LOBBY" | "NIGHT" | "DAY"
  players: map[playerID]*Player   # In-game players
  actions: []Action               # Action queue per phase
}

Player {
  id: string                      # UUID
  name: string                    # Display name
  roleID: int                     # Role from Postgres
  isAlive: bool
  joinedAt: time.Time
}

Action {
  id, playerID, type, targetID, phase, timestamp
}
```

---

## Common Development Tasks

### Task: Implement WebSocket Handler
**File:** `internal/handlers/ws.go`  
**Context:** HTTP upgrade to WebSocket, register client with hub, read/write message pumps  
**Current State:** Empty stub  
**Requirements:**
- Use gorilla/websocket to upgrade connection
- Extract room code & player ID from query params (`/ws?room=CODE&player=ID`)
- Create Client struct, add to hub.register channel
- Read loop: parse JSON messages, dispatch to appropriate game methods
- Write loop: send messages from client.Send channel back to browser
- Handle disconnect: remove from hub.unregister

**Template to request:**
```
Implement the WebSocket upgrade handler in internal/handlers/ws.go.
- Upgrade HTTP GET /ws?room=CODE&player=ID to WebSocket
- Use gorilla/websocket package
- Extract query params, create Client{RoomCode, PlayerID, Send chan Message}
- Register with hub.register <- client
- Implement readPump (parses JSON, calls game methods)
- Implement writePump (sends messages from client.Send channel)
- On disconnect, send to hub.unregister <- client
Use the Hub.BroadcastToRoom() method to notify other players.
```

### Task: Connect Frontend to WebSocket
**File:** `src/lib/ws.js` + `src/routes/room/[code]/+page.svelte`  
**Current State:** `connectWS()` stub; message handler is empty switch  
**Requirements:**
- `connectWS(code, playerId)` opens WebSocket, stores global ws reference
- `handleMessage(msg)` updates Svelte stores based on message.type
- `sendMessage(msg)` sends JSON to backend
- Update stores on events: `player_joined`, `phase_changed`, `action_submitted`, etc.
- Handle reconnect: if WS closes, show "disconnected" UI
- On component mount, call `connectWS()`; on destroy, call `closeWS()`

**Template to request:**
```
Wire up the frontend WebSocket client to the game state stores.
- Implement handleMessage() to dispatch based on msg.type (player_joined, phase_changed, action_submitted, etc.)
- Update store: room (phase, players list), log (new events)
- Add error/disconnect handling: show "reconnecting..." if WS closes
- In src/routes/room/[code]/+page.svelte, call connectWS() in onMount, closeWS() in cleanup
- Ensure sendMessage() works for actions (advance_phase, submit_action, delete_action)
Use the stores.js exports: player, room, isHost, actions, log
```

### Task: Extend HTTP Handlers
**Files:** `internal/handlers/rooms.go`  
**Current Methods:**
- `CreateRoom()` – generates code, calls roomManager.CreateRoom(), returns code + hostID
- `JoinRoom()` – parses JSON body (name), calls roomManager.JoinRoom(), returns playerID + phase
- `UpgradeWebSocket()` – stub

**To extend:**
- Add `StartGame(c)` – host-only endpoint to lock roster and transition to NIGHT phase
- Add `SubmitAction(c)` – player endpoint to queue an action (phase-gated)
- Add `DeleteAction(c)` – host-only endpoint to remove an action from queue

**Template to request:**
```
Extend internal/handlers/rooms.go with three new endpoints:
1. POST /api/rooms/:code/start (host-only) – validate host, start game, assign roles, broadcast via WS
2. POST /api/rooms/:code/actions (player) – validate phase, create Action, add to room.actions, broadcast via WS
3. DELETE /api/rooms/:code/actions/:actionID (host-only) – remove action from queue, broadcast
Use roomHandler.Hub.GetRoomManager() for room state, Hub.BroadcastToRoom() for WS messages
Return error if room not found, player not in room, or permission denied
```

### Task: Add Database Migration
**File:** `internal/db/migrations/NNN_description.up.sql`  
**Process:**
1. Create files: `NNN_description.up.sql` and `NNN_description.down.sql` (e.g., `002_add_game_history.up.sql`)
2. Write up migration (CREATE TABLE, ADD COLUMN, etc.)
3. Write down migration (reverse operation)
4. Run locally: `migrate -path internal/db/migrations -database "postgres://..." up`
5. Test with `sqlc generate` if adding new queries

**Template to request:**
```
Create a migration to [description: add sessions table, add game_history tracking, etc.]
- File: internal/db/migrations/NNN_[description].up.sql
- Include: schema definition with proper types (UUID, TIMESTAMP, JSONB)
- Add .down.sql with DROP TABLE / DROP COLUMN
- Ensure compatibility with existing sqlc queries in internal/db/sql/queries.sql
```

### Task: Add sqlc Query
**File:** `internal/db/sql/queries.sql`  
**Process:**
1. Add SQL comment: `-- name: FunctionName :one` or `:many` or `:exec`
2. Write query with `$1, $2` for parameters
3. Run `sqlc generate` in backend root
4. Use generated code: `querier.FunctionName(ctx, params)`

**Template to request:**
```
Add an sqlc query to internal/db/sql/queries.sql for [purpose: get session by ID, list recent games, etc.]
- Query name: [CamelCase]
- Return type: :one (single row), :many (multiple), :exec (no result)
- Parameters: list as needed (e.g., $1 id, $2 player_name)
After adding, run: sqlc generate
Then use the generated function in your Go code.
```

### Task: Update SvelteKit Page
**Files:** `src/routes/+page.svelte`, `src/routes/room/[code]/+page.svelte`  
**Current State:**
- Landing: create/join forms with basic styling
- Room: player list, host controls (advance phase), event log
- No Skeleton UI integration yet

**To enhance:**
- Add form validation (room code format, player name length)
- Add error feedback (toast/alert on failed join)
- Improve layout with Skeleton UI components (Card, Button, Modal, etc.)
- Wire host panel to actual host-only endpoints

**Template to request:**
```
Enhance src/routes/room/[code]/+page.svelte with:
- Better error handling: catch join/action errors, show toast notifications
- Form validation: check room code and player name before submit
- Host panel improvements: add buttons for AdvancePhase, DeleteAction endpoints
- Skeleton UI integration (optional): use Card, Button, Modal components for better UX
Wire sendMessage() to actual HTTP POST endpoints for StartGame, SubmitAction, DeleteAction
```

---

## Key Code Patterns & Conventions

### WebSocket Message Format (JSON)
```json
{
  "type": "player_joined",
  "data": {
    "playerID": "...",
    "name": "...",
    "role": "..."
  }
}
```

**Common message types:**
- `player_joined` – new player in room
- `player_left` – player disconnected
- `phase_changed` – phase advanced (host action)
- `action_submitted` – player queued action
- `action_deleted` – host removed action
- `roles_assigned` – game started, roles locked
- `game_ended` – game over

### Go Error Handling
```go
if err := someFunc(); err != nil {
  return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
}
```

### Svelte Store Pattern
```js
import { player, room, isHost } from '$lib/stores';

// Subscribe in component
$: console.log($room.phase);  // reactive
room.update(r => ({ ...r, phase: 'DAY' }));  // update
```

### Environment Variables
**Backend:** `internal/` reads `os.Getenv("DATABASE_URL")`, `os.Getenv("PORT")`  
**Frontend:** `import.meta.env.VITE_API_BASE`, `VITE_WS_BASE` (Vite prefix required)

---

## WebSocket Code Commenting Guidelines

All WebSocket code MUST include learning-mode comments explaining the patterns.
This helps developers unfamiliar with WebSockets understand what's happening.

### Why Goroutines for Read/Write Pumps

```go
// WHY GOROUTINE: Each client needs independent read/write loops.
// The read loop blocks waiting for messages from the browser.
// If we didn't use a goroutine, one slow client would block all others.
// Each client gets its own pair of goroutines for reading and writing.
go client.readPump()

// WHY SEPARATE WRITE GOROUTINE: WebSocket connections are NOT thread-safe
// for concurrent writes. If two goroutines try to write at the same time,
// the connection will be corrupted. This goroutine owns all writes to
// this connection, receiving messages via the client.Send channel.
go client.writePump()
```

### Channel Patterns

```go
// CHANNEL PATTERN: Fan-in (many-to-one)
// Multiple client goroutines send to one hub goroutine.
// The hub is the single point of coordination for all rooms.
// This prevents race conditions on the shared rooms map.
hub.register <- client  // Many clients -> one hub

// CHANNEL PATTERN: Fan-out (one-to-many)
// Hub broadcasts to all clients via their individual Send channels.
// Each client has a BUFFERED channel (e.g., size 256) to prevent
// the hub from blocking if one client is slow to read.
client.Send <- message  // One hub -> many clients
```

### Blocking vs Non-Blocking Sends

```go
// BLOCKING SEND: This will wait until the hub processes the registration.
// Safe here because we're in the HTTP handler goroutine, not the hub.
hub.register <- client

// NON-BLOCKING SEND: Use select with default to avoid deadlock.
// If the client's Send buffer is full (client is too slow), we drop
// the message and close the connection rather than blocking the hub.
select {
case client.Send <- msg:
    // Message queued successfully - client will receive it
default:
    // Buffer full - client is too slow, disconnect them
    // This prevents one slow client from blocking all broadcasts
    close(client.Send)
    delete(room, client)
}
```

### Resource Cleanup with Defer

```go
func (c *Client) readPump() {
    // CLEANUP: Always use defer for WebSocket cleanup.
    // This runs when the function returns, even if it panics.
    // Order matters: defer runs in LIFO order (last defer runs first).
    defer func() {
        // 1. Tell the hub we're leaving (removes us from room)
        c.hub.unregister <- c
        // 2. Close the network connection (frees OS resources)
        c.conn.Close()
    }()

    // ... read loop that may exit on error or close ...
}
```

### The Read Loop Explained

```go
func (c *Client) readPump() {
    defer func() { /* cleanup */ }()

    // Set read deadline - if no message received within this time,
    // the read will return an error. This detects dead connections.
    c.conn.SetReadDeadline(time.Now().Add(pongWait))

    // When we receive a "pong" from the browser, reset the deadline.
    // This is the WebSocket ping/pong keepalive mechanism.
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        // ReadMessage blocks until a message arrives or an error occurs.
        // Errors include: connection closed, read deadline exceeded, etc.
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            // Any error means the connection is dead - exit the loop.
            // The defer will clean up.
            break
        }

        // Parse and handle the message...
    }
}
```

### The Write Loop Explained

```go
func (c *Client) writePump() {
    // Ticker sends a ping to the browser every pingPeriod.
    // If the browser doesn't respond with pong, the read loop
    // will timeout and close the connection.
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            // Set write deadline - if write takes too long, fail.
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))

            if !ok {
                // Channel was closed by hub (we're being kicked out).
                // Send a close message to browser and exit.
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            // Write the message to the browser.
            err := c.conn.WriteJSON(message)
            if err != nil {
                return  // Connection dead, exit
            }

        case <-ticker.C:
            // Time to send a ping to keep the connection alive.
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return  // Connection dead, exit
            }
        }
    }
}
```

---

## Testing Strategy

### REQUIRED: Run and verify tests before any code change or commit

**You MUST always run the project’s full test suite before implementing or committing new features, refactoring, or bugfixes. Tests MUST pass. If tests fail, you MUST NOT proceed—stop and fix the test failures first. This applies to all code (backend, frontend, DB migrations, etc.).**

- For any backend database-dependent tests, ensure DB migrations have been applied and TEST_DATABASE_URL is present (see DEV_DB_SETUP.md).
- Use `godotenv` to auto-load environment variables in Go test environments, but remember CLI/migration tools need TEST_DATABASE_URL exported in the shell.
- Do NOT proceed with new tickets until you have a clean test run and all relevant tests pass.

### Priority Order

Test in this order to build confidence incrementally:

1. **Hub unit tests** - Pure Go, no network, no WebSocket
2. **Room unit tests** - Pure Go, no network
3. **WebSocket handler tests** - Mock connections with gorilla test helpers

### Hub Unit Tests (No Network Needed)

```go
// internal/game/hub_test.go

func TestHub_RegisterClient(t *testing.T) {
    hub := NewHub()
    go hub.Run()  // Start the hub's event loop in background

    // Create a test client - no real WebSocket connection needed!
    // We just need the Send channel to exist.
    client := &Client{
        RoomCode: "TEST01",
        PlayerID: "player-123",
        Send:     make(chan Message, 10),  // Buffered channel
    }

    // Register the client by sending to the register channel
    hub.register <- client

    // Give the hub goroutine time to process
    // In production code, you might use sync primitives instead
    time.Sleep(10 * time.Millisecond)

    // Verify the client is now in the room
    hub.mu.Lock()
    clients, exists := hub.rooms["TEST01"]
    hub.mu.Unlock()

    if !exists {
        t.Fatal("room should exist after client registration")
    }
    if len(clients) != 1 {
        t.Fatalf("expected 1 client in room, got %d", len(clients))
    }
}

func TestHub_BroadcastToRoom(t *testing.T) {
    hub := NewHub()
    go hub.Run()

    // Create two clients in the same room
    client1 := &Client{RoomCode: "ROOM1", PlayerID: "p1", Send: make(chan Message, 10)}
    client2 := &Client{RoomCode: "ROOM1", PlayerID: "p2", Send: make(chan Message, 10)}
    // And one client in a different room
    client3 := &Client{RoomCode: "ROOM2", PlayerID: "p3", Send: make(chan Message, 10)}

    hub.register <- client1
    hub.register <- client2
    hub.register <- client3
    time.Sleep(10 * time.Millisecond)

    // Broadcast to ROOM1 only
    hub.BroadcastToRoom("ROOM1", Message{Type: "test", Data: nil})

    // client1 and client2 should receive the message
    select {
    case msg := <-client1.Send:
        if msg.Type != "test" {
            t.Errorf("client1 got wrong message type: %s", msg.Type)
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("client1 didn't receive broadcast")
    }

    select {
    case msg := <-client2.Send:
        if msg.Type != "test" {
            t.Errorf("client2 got wrong message type: %s", msg.Type)
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("client2 didn't receive broadcast")
    }

    // client3 should NOT receive the message (different room)
    select {
    case <-client3.Send:
        t.Error("client3 should not have received broadcast for ROOM1")
    case <-time.After(50 * time.Millisecond):
        // Good - no message received
    }
}
```

### Room Unit Tests

```go
// internal/game/room_test.go

func TestRoomManager_CreateAndJoin(t *testing.T) {
    rm := NewRoomManager()

    // Create a room
    code := rm.CreateRoom("host-123")
    if code == "" {
        t.Fatal("CreateRoom should return a room code")
    }

    // Verify room exists with correct initial state
    room, err := rm.GetRoom(code)
    if err != nil {
        t.Fatalf("GetRoom failed: %v", err)
    }
    if room.HostID != "host-123" {
        t.Errorf("expected HostID 'host-123', got '%s'", room.HostID)
    }
    if room.Phase != "LOBBY" {
        t.Errorf("expected Phase 'LOBBY', got '%s'", room.Phase)
    }

    // Join a player
    playerID, err := rm.JoinRoom(code, "Alice")
    if err != nil {
        t.Fatalf("JoinRoom failed: %v", err)
    }
    if playerID == "" {
        t.Fatal("JoinRoom should return a player ID")
    }

    // Verify player is in the room
    room, _ = rm.GetRoom(code)
    if len(room.Players) != 1 {
        t.Errorf("expected 1 player, got %d", len(room.Players))
    }
}

func TestRoomManager_AdvancePhase(t *testing.T) {
    rm := NewRoomManager()
    code := rm.CreateRoom("host-123")

    // LOBBY -> NIGHT
    err := rm.AdvancePhase(code)
    if err != nil {
        t.Fatalf("AdvancePhase failed: %v", err)
    }

    room, _ := rm.GetRoom(code)
    if room.Phase != "NIGHT" {
        t.Errorf("expected NIGHT, got %s", room.Phase)
    }

    // NIGHT -> DAY
    rm.AdvancePhase(code)
    room, _ = rm.GetRoom(code)
    if room.Phase != "DAY" {
        t.Errorf("expected DAY, got %s", room.Phase)
    }

    // DAY -> NIGHT (cycles back)
    rm.AdvancePhase(code)
    room, _ = rm.GetRoom(code)
    if room.Phase != "NIGHT" {
        t.Errorf("expected NIGHT, got %s", room.Phase)
    }
}
```

### WebSocket Handler Tests (With Mock Connection)

```go
// internal/handlers/ws_test.go

import (
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gorilla/websocket"
    "github.com/labstack/echo/v4"
)

func TestWebSocket_Upgrade(t *testing.T) {
    // Create the hub and handler
    hub := game.NewHub()
    go hub.Run()
    handler := NewRoomHandler(hub)

    // Create Echo server with WebSocket route
    e := echo.New()
    e.GET("/ws", handler.UpgradeWebSocket)

    // Start a test HTTP server
    server := httptest.NewServer(e)
    defer server.Close()

    // Convert HTTP URL to WebSocket URL
    // http://127.0.0.1:xxxxx -> ws://127.0.0.1:xxxxx
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") +
        "/ws?room=TESTROOM&player=player-1"

    // Connect with a WebSocket client
    conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("WebSocket dial failed: %v", err)
    }
    defer conn.Close()

    if resp.StatusCode != 101 {
        t.Errorf("expected status 101 Switching Protocols, got %d", resp.StatusCode)
    }

    // Optionally: send a message and verify response
    // conn.WriteJSON(...)
    // conn.ReadJSON(...)
}

func TestWebSocket_MissingParams(t *testing.T) {
    hub := game.NewHub()
    go hub.Run()
    handler := NewRoomHandler(hub)

    e := echo.New()
    e.GET("/ws", handler.UpgradeWebSocket)

    server := httptest.NewServer(e)
    defer server.Close()

    // Try to connect without room parameter
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?player=p1"

    _, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err == nil {
        t.Error("expected error when room param is missing")
    }
}
```

### Running Tests

```bash
cd backend

# Run all tests with verbose output
go test ./... -v

# Run only game package tests
go test ./internal/game/... -v

# Run only handler tests
go test ./internal/handlers/... -v

# Run tests with coverage report
go test ./... -cover

# Run a specific test by name
go test ./internal/game/... -run TestHub_RegisterClient -v
```

---

## Testing Locally

### Full Stack Test (Do This Before Coding)
```bash
# Terminal 1: Postgres
docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres:15
docker exec postgres createdb -U postgres social_deduction

# Terminal 2: Migrations
cd backend
migrate -path internal/db/migrations -database "postgres://postgres:postgres@localhost:5432/social_deduction" up

# Terminal 3: Backend
go run cmd/server/main.go  # should print "Server running on :8080"

# Terminal 4: Frontend
cd frontend
bun --bun run dev  # should print "http://localhost:5173"

# Browser: Open two windows of http://localhost:5173
# One creates a room, one joins with room code
# Test: Can you see each other? Can host advance phases?
```

### Common Issues & Fixes
| Issue | Fix |
|-------|-----|
| "Connection refused" on backend | Ensure Postgres running; check DATABASE_URL in .env |
| "CORS error" on frontend | Verify Echo CORS config allows GET/POST from 5173 |
| WebSocket won't connect | Check `/ws` endpoint exists; verify room code + player ID in query params |
| Frontend shows "disconnected" | Backend likely crashed; check logs in Terminal 3 |
| sqlc won't generate | Ensure queries.sql has proper comment syntax: `-- name: FunctionName :one` |

---

## Deployment Checklist (When Ready)

- [ ] Update `yourmodule` to actual module path in all Go files
- [ ] Run `go mod tidy` in backend
- [ ] Test locally with multiple browser windows
- [ ] Implement core WebSocket handlers (player_joined, phase_changed, action_submitted)
- [ ] Wire frontend stores to incoming WS messages
- [ ] Create Railway account and install CLI (`railway login`)
- [ ] Create Railway project: `railway init`
- [ ] Link services:
  - Backend: GitHub repo with Dockerfile + go.mod
  - Frontend: SvelteKit app with adapter-node + npm scripts
  - Database: Create Postgres plugin in Railway dashboard
- [ ] Set environment variables:
  - Backend: `DATABASE_URL` (linked auto from Postgres)
  - Frontend: `VITE_API_BASE` (backend URL), `VITE_WS_BASE` (backend WS URL)
- [ ] Deploy: `railway up` or push to GitHub (if connected)
- [ ] Test live: share room code with actual users

---

## Asking for Help (Template)

When requesting AI help on this project, include:

```
Project: Social Deduction Game PoC
Task: [WebSocket handlers / Frontend integration / etc.]
File(s): [internal/handlers/ws.go / src/lib/ws.js / etc.]
Current State: [What exists now / What's missing / What's broken]
Requirements: [What needs to happen / Constraints / Edge cases]
Context: [Go + Echo, WebSocket with gorilla/websocket, hub broadcasts to room]

Code to improve / implement:
[Paste relevant snippet if needed]
```

Example:
```
Project: Social Deduction Game PoC
Task: Implement WebSocket message handling
File: internal/handlers/ws.go
Current State: HTTP upgrade stub; gorilla/websocket imported but not used
Requirements:
- Upgrade GET /ws?room=CODE&player=ID to WebSocket
- Parse incoming JSON messages with format {type, data}
- Dispatch to game.Hub.BroadcastToRoom() based on message type
- Handle graceful disconnect

Context: Hub is a global event loop with register/unregister channels; BroadcastToRoom sends Message struct to all clients in a room.
```

---

## Additional Resources

Backend:
- **golang-migrate:** https://github.com/golang-migrate/migrate (schema management)
- **sqlc:** https://sqlc.dev/ (type-safe queries)
- **gorilla/websocket:** https://github.com/gorilla/websocket (examples in pkg)
- **labstack/echo:** https://echo.labstack.com/docs

Frontend:
- **SvelteKit docs:** https://svelte.dev/docs/kit (load, +layout, stores)
- **Skeleton UI for SvelteKit:** https://www.skeleton.dev/ (components, theming)
- **Bun SvelteKit:** https://bun.com/docs/guides/ecosystem/sveltekit (dev & build)

Deployment:
- **Railway docs:** https://docs.railway.app/guides/sveltekit (deployment)


You are able to use the Svelte MCP server, where you have access to comprehensive Svelte 5 and SvelteKit documentation. Here's how to use the available tools effectively:

## Available MCP Tools:

### 1. list-sections

Use this FIRST to discover all available documentation sections. Returns a structured list with titles, use_cases, and paths.
When asked about Svelte or SvelteKit topics, ALWAYS use this tool at the start of the chat to find relevant sections.

### 2. get-documentation

Retrieves full documentation content for specific sections. Accepts single or multiple sections.
After calling the list-sections tool, you MUST analyze the returned documentation sections (especially the use_cases field) and then use the get-documentation tool to fetch ALL documentation sections that are relevant for the user's task.

### 3. svelte-autofixer

Analyzes Svelte code and returns issues and suggestions.
You MUST use this tool whenever writing Svelte code before sending it to the user. Keep calling it until no issues or suggestions are returned.

### 4. playground-link

Generates a Svelte Playground link with the provided code.
After completing the code, ask the user if they want a playground link. Only call this tool after user confirmation and NEVER if code was written to files in their project.
