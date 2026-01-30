---
name: error-handling-patterns
description: Patterns for consistent error reporting and logging in Go and SvelteKit, including user-facing notifications and backend structured logs.
license: MIT
compatibility: opencode
metadata:
  audience: all
  pattern: error-handling
---

## What I do
- Illustrate structured error propagation in Go (context, HTTP error returns, log best practices)
- Provide SvelteKit/UI-side notification patterns for user-visible errors/toasts
- Checklist for full-stack traceability in multiplayer/websocket flows

## When to use me
Use when designing new features, reviewing code (esp. multiplayer/session flows), or debugging production issues. Prevents silent/hidden failure modes across stack.

## Example Usage
Useful in `backend/internal/handlers/`, `frontend/src/routes/`, and in SvelteKit UI logic that needs robust user feedback.