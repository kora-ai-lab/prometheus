# Prometheus Ghost Shell Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the Tauri-based Ghost Shell UI — a frameless, glassmorphic desktop app with Alt+Space Omnibox, execution modals, and real-time SSE streaming from the Go Core Service.

**Tech Stack:** Tauri v2, Rust (backend), React + TypeScript (frontend), TailwindCSS

---

### Task 1: Tauri Project Scaffolding
**Files:**
- Create: `ghost-shell/` directory with full Tauri project

- [ ] **Step 1.1: Initialize Tauri v2 project**
  - `npm create tauri-app@latest ghost-shell -- --template react-ts`
  - Configure `tauri.conf.json`:
    - Window: frameless, transparent, width 480, height 60 (Omnibox size)
    - Dev path: `../dist` (for Vite build)
    - Bundle identifier: `com.kora-ai.prometheus-shell`

- [ ] **Step 1.2: Configure Tauri permissions**
  - Add `global-shortcut` plugin for Alt+Space
  - Add `shell` plugin for opening files/URLs
  - Add `notification` plugin for task completion alerts

- [ ] **Step 1.3: Setup frontend tooling**
  - Install TailwindCSS
  - Configure Vite for React + TypeScript
  - Add Framer Motion for animations

- [ ] **Step 1.4: Commit**
  `git commit -m "feat: scaffold Tauri v2 project for Ghost Shell"`

### Task 2: Core Service Client
**Files:**
- Create: `ghost-shell/src/lib/api.ts`
- Create: `ghost-shell/src/lib/types.ts`

- [ ] **Step 2.1: Define TypeScript types**
  ```typescript
  interface Task {
    id: string
    goal: string
    status: 'running' | 'blocked' | 'done' | 'failed' | 'cancelled'
    progress: string
    result: string
    error: string
    createdAt: string
    updatedAt: string
  }

  interface HealthResponse {
    status: string
    version: string
    uptime: string
  }
  ```

- [ ] **Step 2.2: Implement API client**
  - `execute(goal: string) => Promise<string>` (POST /api/execute, returns task_id)
  - `getTask(id: string) => Promise<Task>` (GET /api/tasks/{id})
  - `cancelTask(id: string) => Promise<void>` (DELETE /api/tasks/{id})
  - `getHealth() => Promise<HealthResponse>` (GET /api/health)
  - `streamTask(id: string, onEvent: (task: Task) => void) => AbortController` (SSE)
  - Token management: read from `%LOCALAPPDATA%\Prometheus\token.txt`

- [ ] **Step 2.3: Service health monitoring**
  - Poll `/api/health` every 5 seconds
  - Show "Core Service Offline" banner when unreachable
  - Auto-reconnect when service comes back

- [ ] **Step 2.4: Commit**
  `git commit -m "feat: implement core service API client with SSE streaming"`

### Task 3: Omnibox UI Component
**Files:**
- Create: `ghost-shell/src/components/Omnibox.tsx`
- Create: `ghost-shell/src/components/Omnibox.css`

- [ ] **Step 3.1: Build Omnibox shell**
  - Frameless, pill-shaped input bar
  - Centered at bottom of screen
  - Glassmorphism: `backdrop-filter: blur(20px)`, semi-transparent background
  - Gradient border glow (cyan → purple)
  - Input placeholder: "Ask Agent..."
  - Icons: Search (left), Mic, Clip, Expand (right)

- [ ] **Step 3.2: Implement animations**
  - Slide-up on `Alt+Space` (Framer Motion)
  - Fade-out on `Esc`
  - Glow pulse on focus

- [ ] **Step 3.3: Wire to Tauri global shortcut**
  - Register `Alt+Space` via `@tauri-apps/plugin-global-shortcut`
  - Show/hide Omnibox window
  - Auto-focus input on show

- [ ] **Step 3.4: Wire Enter key to API**
  - On Enter: call `execute(goal)`, transition to ExecutionModal
  - Show loading spinner while waiting for task_id

- [ ] **Step 3.5: Commit**
  `git commit -m "feat: build glassmorphic Omnibox with Alt+Space trigger"`

### Task 4: Execution Modal
**Files:**
- Create: `ghost-shell/src/components/ExecutionModal.tsx`
- Create: `ghost-shell/src/components/ResultModal.tsx`

- [ ] **Step 4.1: Thinking/Analyzing state**
  - Centered modal with gradient header (cyan → purple)
  - Animated text showing current progress from SSE stream
  - Thin gradient progress bar
  - "Cancel" button (calls DELETE /api/tasks/{id})
  - States: "Initializing...", "Thinking...", "Executing: ...", "Waiting for input..."

- [ ] **Step 4.2: Result state**
  - Content area with markdown rendering
  - Action buttons: "Copy", "Open File" (if path detected), "Run Again"
  - Smooth transition from thinking → result

- [ ] **Step 4.3: Blocked state**
  - Modal with question from agent
  - Input field for user response
  - Resume button

- [ ] **Step 4.4: Error state**
  - Red-tinted modal
  - Error message display
  - "Retry" button

- [ ] **Step 4.5: Commit**
  `git commit -m "feat: build execution modal with SSE-driven progress"`

### Task 5: System Tray & Settings
**Files:**
- Create: `ghost-shell/src-tauri/src/tray.rs`
- Create: `ghost-shell/src/components/Settings.tsx`

- [ ] **Step 5.1: System tray integration**
  - Custom tray icon
  - Right-click menu: "Open Omnibox", "Settings", "Quit"
  - Double-click tray icon → show Omnibox

- [ ] **Step 5.2: Settings panel**
  - Core service URL configuration (default: http://localhost:8080)
  - Theme toggle (Glass/Dark/Light)
  - Notification preferences
  - API token manual override

- [ ] **Step 5.3: Command history**
  - Store recent goals in localStorage
  - Quick-access dropdown in Omnibox (arrow down to show history)

- [ ] **Step 5.4: Commit**
  `git commit -m "feat: add system tray, settings, and command history"`

### Task 6: Build & Packaging
- [ ] **Step 6.1: Configure Tauri bundling**
  - MSI installer output
  - Icon assets
  - Version bump to 0.1.0

- [ ] **Step 6.2: Build test**
  - `npm run tauri build`
  - Verify MSI output
  - Test installation on clean Windows environment

- [ ] **Step 6.3: Final Commit**
  `git commit -m "feat: configure Tauri bundling and MSI packaging"`

---

## Dependencies Graph
```
Task 1 (Scaffolding) → Task 2 (API Client) → Task 3 (Omnibox) → Task 4 (Modal)
                                                            ↘ Task 5 (Tray/Settings)
Task 5 → Task 6 (Build & Package)
Task 4 → Task 6
```

## Design Tokens
```css
--bg-deep: #0f0f23
--cyan: #00d4ff
--purple: #7b2ff7
--glass-bg: rgba(15, 15, 35, 0.7)
--glass-border: rgba(255, 255, 255, 0.1)
--gradient: linear-gradient(135deg, #00d4ff, #7b2ff7)
```

## Risk Mitigation
- **Tauri v2 compatibility:** Use latest stable plugins
- **SSE reconnection:** Implement exponential backoff
- **Glassmorphism performance:** Use hardware-accelerated CSS transforms
- **Token security:** Store in OS keychain if available, fallback to encrypted file
