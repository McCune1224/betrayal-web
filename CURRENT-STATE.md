# Current Project State

> **Last updated:** January 2025
> 
> This document provides ground truth about what **EXISTS** vs what is **PLANNED**.
> Always check here before assuming code exists. Update this when implementation changes.

---

## Status Legend

| Status | Meaning |
|--------|---------|
| [DONE] | Fully implemented and working |
| [PARTIAL] | File exists, some functionality works |
| [STUB] | File exists, empty or placeholder only |
| [MISSING] | Referenced in AGENTS.md but does not exist |

---

## Backend State

### Entry Point & Server

| File | Status | Notes |
|------|--------|-------|
| `cmd/server/main.go` | [PARTIAL] | Server starts, routes registered, WS endpoint not functional |
| `backend/Dockerfile` | [DONE] | Multi-stage Go build for Railway |
| `backend/railway.toml` | [DONE] | Railway deployment config |
| `backend/.env.example` | [DONE] | DATABASE_URL and PORT |
| `backend/go.mod` | [DONE] | Module path: betrayal-web, all deps present |

### Internal Packages

| File | Status | Notes |
|------|--------|-------|
| `internal/db.go` | [DONE] | Database connection with pgx |
| `internal/models/types.go` | [DONE] | Room, Player, Action, Role structs defined |
| `internal/game/room.go` | [PARTIAL] | CreateRoom, JoinRoom, AdvancePhase work. Missing: StartGame, LeaveRoom, SubmitAction, DeleteAction |
| `internal/game/hub.go` | [PARTIAL] | Hub struct, register/unregister channels, BroadcastToRoom exist. Not integrated with actual WebSocket connections |
| `internal/game/errors.go` | [PARTIAL] | Only 2 errors defined (ErrRoomNotFound, ErrInvalidRoom). Missing many error types |
| `internal/handlers/rooms.go` | [PARTIAL] | CreateRoom, JoinRoom handlers work. Missing: StartGame, SubmitAction, DeleteAction endpoints |
| `internal/handlers/ws.go` | [STUB] | **Empty file** - no WebSocket implementation |

### Database Layer

| File | Status | Notes |
|------|--------|-------|
| `internal/db/migrations/001_init_schema.up.sql` | [DONE] | roles + sessions tables |
| `internal/db/migrations/001_init_schema.down.sql` | [DONE] | Rollback migration |
| `internal/db/sql/queries.sql` | [PARTIAL] | Only 2 queries: GetRoleByID, ListRoles |
| `internal/db/sqlc/` | [DONE] | Generated: db.go, models.go, queries.sql.go |
| `sqlc.yaml` | [DONE] | sqlc configuration exists |

### Dependencies (go.mod)

| Dependency | Status | Notes |
|------------|--------|-------|
| `github.com/labstack/echo/v4` | [DONE] | HTTP framework |
| `github.com/jackc/pgx/v5` | [DONE] | Postgres driver |
| `github.com/joho/godotenv` | [DONE] | Env file loading |
| `github.com/google/uuid` | [DONE] | UUID generation |
| `github.com/gorilla/websocket` | [DONE] | WebSocket library added |

### Tests

| Directory | Status | Notes |
|-----------|--------|-------|
| `internal/game/*_test.go` | [DONE] | 32 tests: hub_test.go (15), room_test.go (17) |
| `internal/handlers/*_test.go` | [MISSING] | **No handler tests exist** |

---

## Frontend State

### Routes

| File | Status | Notes |
|------|--------|-------|
| `src/routes/+page.svelte` | [STUB] | **Default SvelteKit welcome page** - no game UI |
| `src/routes/+layout.svelte` | [DONE] | Root layout with CSS imports |
| `src/routes/room/[code]/+page.svelte` | [MISSING] | **Game room page doesn't exist** |

### Libraries

| File | Status | Notes |
|------|--------|-------|
| `src/lib/stores.js` | [MISSING] | **No Svelte stores defined** |
| `src/lib/api.js` | [MISSING] | **No HTTP client** |
| `src/lib/ws.js` | [MISSING] | **No WebSocket client** |
| `src/lib/index.ts` | [STUB] | Empty barrel file |

### Configuration

| File | Status | Notes |
|------|--------|-------|
| `frontend/.env.example` | [DONE] | VITE_API_BASE and VITE_WS_BASE defined |
| `frontend/Dockerfile` | [MISSING] | **Cannot deploy frontend** |
| `svelte.config.js` | [DONE] | adapter-node configured for Railway |
| `tailwind.config.js` | [DONE] | Tailwind v4 configured |

---

## Project-Level Files

| File | Status | Notes |
|------|--------|-------|
| `AGENTS.md` | [DONE] | Target architecture (aspirational) |
| `CURRENT-STATE.md` | [DONE] | This file (ground truth) |
| `MVP-GOALS.md` | [DONE] | Checklist tracker |
| `project-brief.md` | [DONE] | Requirements/specification |
| `README.md` | [MISSING] | No project readme |
| `SETUP_GUIDE.md` | [MISSING] | No setup instructions |

---

## Blocking Issues

These must be resolved before significant progress can be made:

1. **No frontend lib files**
   - Cannot connect frontend to backend
   - Fix: Create stores.js, api.js, ws.js

2. **WebSocket handler not implemented**
   - ws.go is still an empty stub
   - Fix: Implement WebSocket upgrade and read/write pumps

---

## Next Actions

See **MVP-GOALS.md** for prioritized checklist of what to implement next.

---

## How to Update This File

When you implement a feature:
1. Change the status label ([MISSING] -> [STUB] -> [PARTIAL] -> [DONE])
2. Update the Notes column
3. Update the "Last updated" date at the top
4. Remove from "Blocking Issues" if resolved
