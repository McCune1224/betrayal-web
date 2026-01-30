---
name: websocket-handler-go
description: Step-by-step guide for implementing Go WebSocket handlers with gorilla/websocket, including channel patterns, resource cleanup, and thorough learning-mode comments.
license: MIT
compatibility: opencode
metadata:
  audience: backend-dev
  pattern: websocket-hub
---

## What I do
- Show idiomatic Go WebSocket handler scaffolding with AGENTS.md-based learning-mode comments
- Document goroutines, fan-in/fan-out channel and buffer patterns, read/write lifecycle
- Checklist for Client registration, unregister, error handling, resource cleanup
- Example integration test/skeleton (mock websocket)

## When to use me
Use this skill when implementing, refactoring, or reviewing multiplayer WebSocket network code in Go, onboarding backend devs, or enforcing robust connection management/documentation.

## Example Usage
Refer to this when creating or maintaining files like `backend/internal/handlers/ws.go`, or when scaffolding a new game hub/event loop. It helps ensure lifecycle edge cases (disconnects, slow clients, cleanup) are not missed.