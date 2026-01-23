# Betrayal - Social Deduction Game (PoC) – Project Brief

## Overview

A turn-based, multiplayer social deduction game inspired by Blood on the Clock Tower, but with every player competing individually (battle royale). Players have roles with abilities and perks. A designated host advances game phases and manages actions. Real-time updates via WebSockets.

## Tech Stack

**Backend:**
- Language: Go
- Framework: Echo (HTTP API + WebSocket support)
- WebSockets: gorilla/websocket + Echo's cookbook pattern
- Database: Postgres (Railway-hosted)
- Migrations: golang-migrate
- Data access: sqlc (code generation) + pgx (driver) + optional sqlx (ad-hoc queries)
- Messages: JSON with `type` field

**Frontend:**
- Framework: SvelteKit
- Hosting: Railway (separate service from backend)
- WebSocket client: browser native WebSocket API

**Infrastructure:**
- Hosting: Railway (Go backend + SvelteKit frontend + Postgres)
- State: In-memory room maps (no persistence for rooms)
- Auth: Anonymous session IDs (player name + auto-generated ID for reconnect)

## Game Model

### Core Mechanics

- **Unbounded player count** per room (practically limited by single-node server)
- **Room creation:** Manual room codes (e.g., "ABC123")
- **Host:** Out-of-game (storyteller-style), cannot win, can control game flow
- **Roles:** Pre-configured in Postgres, no duplicates per room
- **Phases:** `LOBBY` → `NIGHT` → `DAY` → repeat (host-driven advancement)
- **Actions:** Players submit phase-gated actions (e.g., "kill player X", "protect player Y")
- **Action resolution:** Host can approve or delete actions (simple rollback)
- **Reconnect:** Player can rejoin via same room code + player name; if seat is free, they reclaim it
- **Host disconnect:** Room freezes until host returns (no auto-promotion)

### Data Model (Sketch)

**Roles table** (persisted, seeded manually):
- `id`, `name`, `alignment`, `description`, `perks_json`, `starting_items_json`

**In-memory room state:**
- Room code, host player ID, phase, player roster (UUID + name + role + is_alive, etc.), action queue

**Sessions** (optional for later):
- Player ID, player name, last join timestamp (for reconnect tracking)

## MVP Features

### Core Gameplay (Backend + Frontend)

- [ ] Create room (HTTP `POST /api/rooms`): returns room code
- [ ] Join room (HTTP `POST /api/rooms/{code}/join`): returns room state + player ID
- [ ] WebSocket connection (`/ws?room=CODE&player=ID`):
  - Client subscribes to room updates
  - Server broadcasts:
    - Player joined / left
    - Phase changed
    - Action submitted / deleted
    - Player reconnected
- [ ] Host-only controls:
  - Start game (locks roster, assigns available roles, transitions to first phase)
  - Advance phase
  - Delete/rollback an action
- [ ] Player actions:
  - Submit action during allowed phases
  - View current phase and allowed actions
  - See private role info (hidden from others)
- [ ] Host private panel:
  - List all players + their roles (visible to host only)
  - Action log with player + target + action type

### Database & Schema

- [ ] Postgres Roles table with basic fields
- [ ] golang-migrate setup with initial schema
- [ ] sqlc queries for roles, sessions (optional)

### Frontend (SvelteKit)

- [ ] Landing page: create / join room forms
- [ ] Lobby view: player roster, host-only start button
- [ ] In-game view:
  - Current phase display
  - Role card (hidden if not started)
  - Action panel (phase-dependent)
  - Event log
  - Host panel (if host): player list + roles, action history
- [ ] Basic styling (no heavy design yet, functional UI)

### Hosting (Railway)

- [ ] Go backend deployed to Railway
- [ ] SvelteKit frontend deployed to Railway
- [ ] Postgres instance linked to Go app (DATABASE_URL)
- [ ] CORS configured on Echo for cross-origin WebSocket + HTTP

## Non-MVP / Future (Out of Scope for PoC)

- Admin UI for creating / editing roles
- Spectator mode
- Persistent room history / game statistics
- Voice / video integration
- Cosmetics, progression, unlocks
- Protobuf message encoding (upgrade from JSON)
- Horizontal scaling / multi-node support

## Development Milestones (Rough)

**Week 1: Backend foundations**
- [ ] Go project structure, Echo setup, Postgres migrations
- [ ] WebSocket hub (room registry, broadcast helpers)
- [ ] Basic HTTP endpoints (create room, join room)
- [ ] Schema: roles table, roles queries (sqlc)

**Week 2: Game logic + WebSocket messages**
- [ ] Room state machine (phase transitions, role assignment)
- [ ] Action submission / deletion
- [ ] WebSocket message handlers (join, phase change, action submit, etc.)
- [ ] Host-only endpoint guards

**Week 3: Frontend + integration**
- [ ] SvelteKit landing + room join flow
- [ ] In-game view + WebSocket client
- [ ] Host panel + player roster
- [ ] Basic event log / action history

**Week 4: Polish + deploy**
- [ ] Bug fixes, edge cases (host disconnect, etc.)
- [ ] Deploy to Railway
- [ ] Manual testing with real players / browsers

## Key Decisions Locked

- JSON over WebSockets (not Protobuf, for now)
- In-memory room state (no persistence for rooms across server restarts)
- Host is out-of-game and cannot win
- Simple "delete and re-enter" action rollback (no complex undo history)
- Roles defined in DB, seeded manually (no admin UI for PoC)
- Basic reconnect via player name + room code
- Host disconnect freezes room (no auto-promotion)

## Success Criteria for PoC

1. Players can create a room and share a code
2. Players can join and see each other in real-time
3. Host can start game, assign roles, and advance phases
4. Players can submit actions in the correct phases
5. Host can see all roles and manage action queue
6. Reconnect works: player rejoins with same name/ID
7. Deployed on Railway with live WebSocket support
8. No crashes on typical flows (8–16 players, 2–3 phase cycles)
