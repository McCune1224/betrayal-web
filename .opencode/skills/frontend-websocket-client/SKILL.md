---
name: frontend-websocket-client
description: Blueprint for SvelteKit/WebSocket integration, including reconnect logic, error display, and connection-state-driven Svelte store patterns.
license: MIT
compatibility: opencode
metadata:
  audience: frontend-dev
  pattern: ws-client
---

## What I do
- Show robust connectWS/sendMessage/handleMessage patterns in SvelteKit using native WebSocket API and stores
- Provide approaches for error handling, reconnection UI (disconnected/reconnected banners), and decoupled store updates
- Illustrate the correct cleanup of websockets (on component destroy)

## When to use me
Use when enabling any real-time front-end interaction (e.g., room/player/game phase updates) powered by WebSockets, or reviewing new event-driven UI.

## Example Usage
Reference when working on `frontend/src/lib/ws.js`, `/room/[code]/+page.svelte`, or UI components that depend on live state updates.