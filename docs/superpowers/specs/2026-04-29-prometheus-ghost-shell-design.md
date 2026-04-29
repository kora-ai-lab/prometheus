# Design Spec: Prometheus Ghost Shell (Frontend UI)
Date: 2026-04-29
Status: Proposed

## 1. Visual Identity & Philosophy
The UI follows a "Ghost" aesthetic: minimal, floating, and high-contrast. It avoids traditional window frames in favor of floating elements with glassmorphism effects.

### Design Tokens
- **Backgrounds:** Glassmorphism (semi-transparent blur), Deep Navy (#0f0f23), Linear Gradients (Cyan #00d4ff $\rightarrow$ Purple #7b2ff7).
- **Typography:** San Francisco / Segoe UI, Clean sans-serif.
- **Animations:** Smooth transitions (Ease-in-out), slide-up for Omnibox, fade-in for modals.

## 2. Component Specifications

### 2.1 The Omnibox (Input Bar)
- **Trigger:** Global shortcut `Alt + Space`.
- **Visuals:** 
    - Floating centered bar at the bottom of the screen.
    - Rounded corners (pill-shaped).
    - Subtle outer glow and backdrop blur.
- **Interactions:**
    - `Enter` $\rightarrow$ Trigger core execution.
    - `Esc` $\rightarrow$ Hide bar.
- **Elements:**
    - Search icon / Avatar.
    - Input field: "Ask Agent...".
    - Action buttons: Voice input (Mic icon), Attachment (Clip icon), Fullscreen (Expand icon).

### 2.2 Execution Modals (The "Ghost" Windows)
Instead of a terminal, the agent communicates through elegant modal overlays.

- **Thinking/Analyzing State:**
    - Centered modal with a gradient header.
    - Dynamic text: "Analyzing document...", "Searching the web...".
    - Progress bar: Thin, animated gradient line.
    - "Cancel" button.
- **Result State:**
    - Content area displaying the agent's a-priori result.
    - Contextual Action Buttons: "Open File", "Copy", "Run Again".
    - Bottom input for follow-up questions.

### 2.3 Control Center & History
- **System Tray:** Icon in the Windows tray with a right-click menu.
- **Settings Panel:** Modern toggle-based interface for:
    - Context Sources (Desktop, Browser, Local Files).
    - Personalization (Light/Dark/Glass theme).
- **Command History:** List of past goals with status tags (Completed, Failed, In Progress).

## 3. Technical Stack (Proposed)
- **Framework:** Tauri (Rust + React/Next.js).
- **Why?**
    - **Ultra-lightweight:** Much smaller binaries than Electron (matches the <<225MB goal).
    - **Native Windows API:** Better control over global shortcuts, tray icons, and frameless windows.
    - **Web-based UI:** Allows us to implement the exact CSS gradients and glassmorphism from the design images.

## 4. Interaction Flow
`Alt + Space` $\rightarrow$ `Type Goal` $\rightarrow$ `Enter` $\rightarrow$ `Show Analysis Modal` $\rightarrow$ `Show Result Modal`
