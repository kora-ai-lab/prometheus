# Auto-Launch Browser Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Automatically launch Chrome or Edge with remote debugging enabled if not already running, eliminating the need for manual browser startup.

**Architecture:** 
1. Create a `launcher` module to handle browser discovery (finding the executable) and process management.
2. Implement a check to see if port 9222 is active.
3. Integrate this check into `WebDriver::connect` to ensure the browser is up before attempting the WebSocket connection.

**Tech Stack:** Rust, `tokio::process`, `anyhow`, `std::net::TcpStream`.

---

### Task 1: Browser Discovery Module

**Files:**
- Create: `src/web/launcher.rs`

- [ ] **Step 1: Implement `find_browser_executable()`**
  - Search common Windows paths for `chrome.exe` and `msedge.exe`.
  - Return the first one found as a `PathBuf`.

- [ ] **Step 2: Implement `is_port_open(port: u16)`**
  - Attempt a `TcpStream::connect` to `localhost:port`.
  - Return `true` if connection succeeds, `false` otherwise.

- [ ] **Step 3: Implement `launch_browser(path: PathBuf)`**
  - Use `tokio::process::Command` to start the browser.
  - Add argument `--remote-debugging-port=9222`.
  - Use `.spawn()` to run it in the background.

- [ ] **Step 4: Implement `ensure_browser_running()`**
  - Check `is_port_open(9222)`.
  - If closed, find executable and launch it.
  - Wait a few seconds for the browser to initialize.

- [ ] **Step 5: Commit**
  - `git add src/web/launcher.rs`
  - `git commit -m "feat: add browser launcher for auto-start"`

### Task 2: Integration with WebDriver

**Files:**
- Modify: `src/web/cdp.rs`
- Modify: `src/web/mod.rs` (if exists, to declare the new module)

- [ ] **Step 1: Declare the `launcher` module in `src/web/mod.rs` (or `src/main.rs` if not using mod.rs)**

- [ ] **Step 2: Update `WebDriver::connect(port)`**
  - Call `launcher::ensure_browser_running().await?` at the start of the function.

- [ ] **Step 3: Commit**
  - `git add src/web/cdp.rs src/web/mod.rs`
  - `git commit -m "feat: integrate auto-launch into WebDriver connection"`

### Task 3: Verification

- [ ] **Step 1: Run the "Google title" test**
  - Command: `C:\Users\junio\.cargo\bin\cargo.exe run -- "Go to google.com and tell me the title of the page"`
  - Verify: Browser launches automatically and the title is returned.

- [ ] **Step 2: Verify no duplicate browsers**
  - Run the command again.
  - Verify that it uses the existing browser instead of launching a new one.
