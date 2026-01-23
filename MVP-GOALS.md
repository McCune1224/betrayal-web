# MVP Goals & Checklist

> Track progress toward minimum viable product.
> Update checkboxes as work completes.

---

## MVP Definition

A working game where:
- [ ] Multiple players can join a room via browser
- [ ] Host can advance game phases
- [ ] Players see real-time updates via WebSocket
- [ ] Basic role assignment works
- [ ] Game can complete (someone wins)

---

## Phase 1: Backend Foundation

> **Goal:** Testable backend logic, no WebSocket implementation yet

### Dependencies (Do First)

- [x] Add `gorilla/websocket` to go.mod: `go get github.com/gorilla/websocket`
- [x] Run `sqlc generate` in backend directory
- [x] Update module path from `yourmodule` to actual path
- [x] Run `go mod tidy`

### Unit Tests (Priority 1 - Before WebSocket)

Testing is the priority because you're unfamiliar with WebSockets. Write tests first to understand the patterns.

**Hub Tests (`internal/game/hub_test.go`):**
- [x] Test Hub.Run() starts without error
- [x] Test client registration adds to room
- [x] Test client unregistration removes from room
- [x] Test BroadcastToRoom sends to all clients in room
- [x] Test BroadcastToRoom doesn't send to other rooms
- [x] Test empty room cleanup after last client leaves

**Room Tests (`internal/game/room_test.go`):**
- [x] Test CreateRoom generates unique codes
- [x] Test JoinRoom adds player to room
- [x] Test JoinRoom fails for non-existent room
- [x] Test GetRoom returns correct room
- [x] Test DeleteRoom removes room
- [x] Test AdvancePhase cycles LOBBY->NIGHT->DAY->NIGHT
- [x] Test AdvancePhase fails for non-existent room

### HTTP Handlers (Verify Existing)

- [ ] Test `POST /api/rooms` creates room and returns code
- [ ] Test `POST /api/rooms/:code/join` returns playerID + phase
- [ ] Test join fails for non-existent room code
- [ ] Test join fails for invalid request body

---

## Phase 2: WebSocket Implementation

> **Goal:** Real-time communication works with heavily commented code

### Add gorilla/websocket

- [x] `go get github.com/gorilla/websocket`
- [x] Verify import works in ws.go

### Handler Implementation (`internal/handlers/ws.go`)

All code should have **learning-mode comments** explaining goroutines, channels, and patterns.

- [x] Implement HTTP upgrade to WebSocket connection
- [x] Extract room code and player ID from query params
- [x] Create Client struct with Send channel
- [x] Implement `readPump()` goroutine with comments explaining:
  - Why it runs in a separate goroutine
  - How it parses JSON messages
  - What happens on read error
- [x] Implement `writePump()` goroutine with comments explaining:
  - Why writes need their own goroutine
  - The ticker for ping/pong keepalive
  - Non-blocking select pattern
- [x] Implement graceful disconnect handling with defer
- [x] Register client with Hub on connection
- [x] Unregister client with Hub on disconnect

### Hub Integration

- [x] Hub.Run() processes register channel correctly
- [x] Hub.Run() processes unregister channel correctly
- [x] Hub.BroadcastToRoom() sends to all room clients
- [x] Client.Send channel receives broadcast messages

### WebSocket Tests (`internal/handlers/ws_test.go`)

- [x] Test WebSocket upgrade succeeds with valid params
- [x] Test upgrade fails without room code
- [x] Test upgrade fails without player ID
- [x] Test message routing to correct room
- [x] Test client disconnect cleanup
- [x] Test broadcast reaches all room members

---

## Phase 3: Frontend Integration

> **Goal:** Browser can create/join rooms and see updates

### Library Files (Create)

**Stores (`src/lib/stores.js`):**
- [ ] `player` store - { id, name, roleId, isAlive }
- [ ] `room` store - { code, phase, players, hostId }
- [ ] `isHost` derived store - computed from player.id === room.hostId
- [ ] `actions` store - action queue array
- [ ] `log` store - event log array

**API Client (`src/lib/api.js`):**
- [ ] `createRoom()` - POST /api/rooms, returns { code, hostId }
- [ ] `joinRoom(code, name)` - POST /api/rooms/:code/join, returns { playerId, phase }
- [ ] Error handling for failed requests

**WebSocket Client (`src/lib/ws.js`):**
- [ ] `connectWS(code, playerId)` - opens WebSocket connection
- [ ] `closeWS()` - closes connection
- [ ] `sendMessage(msg)` - sends JSON to backend
- [ ] `handleMessage(msg)` - dispatches based on msg.type:
  - `player_joined` → update room.players
  - `player_left` → update room.players
  - `phase_changed` → update room.phase
  - `action_submitted` → update actions
  - `action_deleted` → update actions
  - `roles_assigned` → update player.roleId
- [ ] Disconnect/reconnect handling with UI feedback

### Pages

**Landing Page (`src/routes/+page.svelte`):**
- [ ] Create room form with button
- [ ] Join room form with code input and name input
- [ ] Form validation (code format, name length)
- [ ] Error feedback on failed join
- [ ] Navigate to /room/[code] on success

**Room Page (`src/routes/room/[code]/+page.svelte`):**
- [ ] Connect to WebSocket on mount
- [ ] Disconnect on component destroy
- [ ] Show player list from room store
- [ ] Show current phase
- [ ] Show event log
- [ ] Host controls (visible only to host):
  - Start game button
  - Advance phase button
- [ ] Handle disconnected state

---

## Phase 4: Game Logic

> **Goal:** Actual game mechanics work

### Role System

- [ ] Roles seeded in database
- [ ] `StartGame()` assigns random roles to players
- [ ] Player can see their own role
- [ ] Roles hidden from other players (except host)

### Phase Management

- [ ] LOBBY → NIGHT transition (host starts game)
- [ ] NIGHT → DAY transition
- [ ] DAY → NIGHT transition
- [ ] Broadcast phase_changed on each transition

### Actions

- [ ] `SubmitAction()` - player queues action during correct phase
- [ ] `DeleteAction()` - host removes action from queue
- [ ] Actions resolve on phase change
- [ ] Broadcast action_submitted and action_deleted

### Win Conditions

- [ ] Track player alive status
- [ ] Detect win condition after action resolution
- [ ] Broadcast game_ended with winner

---

## Phase 5: Documentation & Deployment

> **Goal:** Others can run and deploy the project

### Documentation

- [ ] Code comments in ws.go explain all patterns
- [ ] Code comments in hub.go explain channel usage
- [ ] README.md with project overview (optional)

### Deployment Files

- [ ] `frontend/.env.example` exists
- [ ] `frontend/Dockerfile` created
- [ ] Backend Dockerfile verified working
- [ ] Railway deployment tested end-to-end

---

## Testing Cheat Sheet

### Running Backend Tests

```bash
cd backend
go test ./internal/game/... -v        # Hub and Room tests
go test ./internal/handlers/... -v    # HTTP and WebSocket tests
go test ./... -v                       # All tests
```

### Hub Test Pattern (No Network)

```go
func TestHub_RegisterClient(t *testing.T) {
    hub := NewHub()
    go hub.Run()  // Start hub event loop

    // Create a test client (no real WebSocket needed)
    client := &Client{
        RoomCode: "TEST01",
        PlayerID: "player-123",
        Send:     make(chan Message, 10),
    }

    // Register the client
    hub.register <- client

    // Give hub time to process
    time.Sleep(10 * time.Millisecond)

    // Assert client is in the room
    hub.mu.Lock()
    clients, exists := hub.rooms["TEST01"]
    hub.mu.Unlock()

    if !exists {
        t.Fatal("room should exist after registration")
    }
    if len(clients) != 1 {
        t.Fatalf("expected 1 client, got %d", len(clients))
    }
}
```

### Room Test Pattern

```go
func TestRoomManager_CreateRoom(t *testing.T) {
    rm := NewRoomManager()

    code := rm.CreateRoom("host-123")

    if code == "" {
        t.Fatal("expected room code, got empty string")
    }

    room, err := rm.GetRoom(code)
    if err != nil {
        t.Fatalf("GetRoom failed: %v", err)
    }
    if room.HostID != "host-123" {
        t.Errorf("expected hostID 'host-123', got '%s'", room.HostID)
    }
    if room.Phase != "LOBBY" {
        t.Errorf("expected phase 'LOBBY', got '%s'", room.Phase)
    }
}
```

### WebSocket Test Pattern (With Mock)

```go
func TestWS_Upgrade(t *testing.T) {
    // Create test server with WebSocket handler
    hub := NewHub()
    go hub.Run()
    handler := NewRoomHandler(hub)

    e := echo.New()
    e.GET("/ws", handler.UpgradeWebSocket)

    server := httptest.NewServer(e)
    defer server.Close()

    // Convert HTTP URL to WebSocket URL
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + 
             "/ws?room=TEST&player=player-1"

    // Connect WebSocket client
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("WebSocket dial failed: %v", err)
    }
    defer conn.Close()

    // Verify connection works by sending/receiving
    // ...
}
```

---

## Progress Summary

| Phase | Status | Blockers |
|-------|--------|----------|
| Phase 1: Backend Foundation | Complete | None |
| Phase 2: WebSocket Implementation | Complete | None |
| Phase 3: Frontend Integration | Not Started | Need to create frontend files |
| Phase 4: Game Logic | Not Started | Phase 3 incomplete |
| Phase 5: Deployment | Not Started | Phase 4 incomplete |
