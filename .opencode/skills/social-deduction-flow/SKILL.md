---
name: social-deduction-flow
description: Enforcement of social deduction game state, role assignment, phase transitions, and workflow invariants for multi-phase game flows.
license: MIT
compatibility: opencode
metadata:
  audience: all
  workflow: game-logic
---

## What I do
- Document the phase/state transition logic for Betrayal-style games (LOBBY, NIGHT, DAY phases)
- Checklist for assigning roles, validating player action permissions by phase and tracking/validating end-of-game detection
- Examples for how state should be synchronized between Go and Svelte frontends

## When to use me
Use when implementing, reviewing, or refactoring core social deduction game logic, especially flow-critical state/transition handling in Go or SvelteKit.

## Example Usage
Reference when augmenting `backend/internal/game/`, and when wiring up frontend state handling for phase or role presentations.