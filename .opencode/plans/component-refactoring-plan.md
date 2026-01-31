# Component Refactoring Plan: Break Down Giant Pages into Reusable Components

## Current State Analysis

### Problem
Both main pages are "giant blob pages":
- `room/[code]/+page.svelte`: 343 lines
- `+page.svelte` (landing): 242 lines

These are difficult to maintain, test, and reuse. Need to extract reusable components.

---

## Proposed Component Structure

### 1. UI Components (Generic/Shared)
**Location:** `$lib/components/ui/`

| Component | Purpose | Props |
|-----------|---------|-------|
| `Button.svelte` | Reusable button with variants | variant, size, disabled, onclick, children |
| `Input.svelte` | Text input with label | label, value, placeholder, type, error, maxlength |
| `Card.svelte` | Container with border/shadow | children, className |
| `Avatar.svelte` | User avatar with initials | name, size, variant |
| `Alert.svelte` | Error/info messages | type (error/success/info), message |

### 2. Feature Components (Domain-Specific)
**Location:** `$lib/components/`

#### From Room Page (`room/[code]/+page.svelte`):

**ConnectionStatus.svelte** (Lines 107-135)
- Displays connection banner
- Props: `status`, `error`, `onReconnect`

**RoomSidebar.svelte** (Lines 142-254)
- Container for left sidebar
- Props: `roomCode`, `phase`, `players`, `isHost`, `onStartGame`, `onAdvancePhase`

**RoomHeader.svelte** (Lines 146-163)
- Room code and phase display
- Props: `roomCode`, `phase`

**PlayerList.svelte** (Lines 165-220)
- List of all players
- Props: `players`, `currentPlayerId`, `isHost`

**PlayerCard.svelte** (Lines 174-194, 198-217)
- Individual player display
- Props: `name`, `isHost`, `isCurrentPlayer`, `variant`

**HostControls.svelte** (Lines 222-253)
- Host action buttons
- Props: `phase`, `isHost`, `disabled`, `onStartGame`, `onAdvancePhase`

**ChatArea.svelte** (Lines 256-341)
- Main chat container
- Props: `messages`, `currentPlayerId`, `onSendMessage`

**ChatMessage.svelte** (Lines 271-315)
- Individual chat message
- Props: `message`, `isOwnMessage`

**ChatInput.svelte** (Lines 319-340)
- Input field with send button
- Props: `value`, `placeholder`, `disabled`, `onSend`, `onKeydown`

#### From Landing Page (`+page.svelte`):

**LandingHeader.svelte** (Lines 103-113)
- Title and tagline
- Props: `title`, `tagline`

**TabSwitcher.svelte** (Lines 119-139)
- Join/Create tabs
- Props: `activeTab`, `tabs`, `onTabChange`

**JoinForm.svelte** (Lines 150-193)
- Join room form
- Props: `username`, `roomCode`, `error`, `isLoading`, `onSubmit`, `onInput`

**CreateForm.svelte** (Lines 194-233)
- Create room form
- Props: `username`, `error`, `isLoading`, `onSubmit`, `onInput`

---

## Implementation Details

### New File Structure
```
frontend/src/lib/
├── components/
│   ├── ui/
│   │   ├── Button.svelte
│   │   ├── Input.svelte
│   │   ├── Card.svelte
│   │   ├── Avatar.svelte
│   │   └── Alert.svelte
│   ├── ConnectionStatus.svelte
│   ├── RoomSidebar.svelte
│   ├── RoomHeader.svelte
│   ├── PlayerList.svelte
│   ├── PlayerCard.svelte
│   ├── HostControls.svelte
│   ├── ChatArea.svelte
│   ├── ChatMessage.svelte
│   ├── ChatInput.svelte
│   ├── LandingHeader.svelte
│   ├── TabSwitcher.svelte
│   ├── JoinForm.svelte
│   └── CreateForm.svelte
```

### Step-by-Step Migration

**Phase 1: Extract UI Components (Foundation)**
1. Create `Button.svelte` - Base button with variants
2. Create `Input.svelte` - Input with label
3. Create `Card.svelte` - Container component
4. Create `Avatar.svelte` - User avatar
5. Create `Alert.svelte` - Error messages

**Phase 2: Extract Room Page Components**
1. Create `ChatMessage.svelte` - Simplest, standalone
2. Create `ChatInput.svelte` - Uses Button/Input
3. Create `ChatArea.svelte` - Uses ChatMessage/ChatInput
4. Create `PlayerCard.svelte` - Uses Avatar
5. Create `PlayerList.svelte` - Uses PlayerCard
6. Create `RoomHeader.svelte` - Simple display
7. Create `HostControls.svelte` - Uses Button
8. Create `RoomSidebar.svelte` - Combines RoomHeader, PlayerList, HostControls
9. Create `ConnectionStatus.svelte` - Uses Button

**Phase 3: Extract Landing Page Components**
1. Create `LandingHeader.svelte`
2. Create `TabSwitcher.svelte`
3. Create `JoinForm.svelte` - Uses Input, Button, Alert
4. Create `CreateForm.svelte` - Uses Input, Button, Alert

**Phase 4: Refactor Pages**
1. Update `room/[code]/+page.svelte` to use new components
2. Update `+page.svelte` to use new components
3. Delete old inline code

---

## Expected Result

### Before:
- `room/[code]/+page.svelte`: 343 lines
- `+page.svelte`: 242 lines
- Total: 585 lines of dense, mixed logic + markup

### After:
- `room/[code]/+page.svelte`: ~80-100 lines (imports + state + composition)
- `+page.svelte`: ~60-80 lines (imports + state + composition)
- 15-20 small, focused components (~20-50 lines each)
- Clear separation of concerns
- Easy to test, maintain, and reuse

---

## Testing Strategy

1. After each component extraction, verify the UI still works
2. Test components in isolation where possible
3. Ensure WebSocket functionality remains intact
4. Verify form validation still works
5. Check responsive behavior

---

## Notes

- Keep existing stores (`$lib/stores.svelte`) - no changes needed
- Keep existing API/WebSocket logic (`$lib/api.ts`, `$lib/ws.ts`)
- Components should be presentational (receive props, emit events)
- Pages handle state management and pass data down to components
- Use TypeScript interfaces for all component props
