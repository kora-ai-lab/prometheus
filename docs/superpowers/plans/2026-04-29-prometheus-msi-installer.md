# Prometheus MSI Installer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a unified Windows installer that bundles the Go Core Service and Ghost Shell UI, registers the service, and handles clean uninstallation.

**Approach:** Enhanced Inno Setup (existing `scripts/windows/prometheus.iss`) — simpler than WiX, supports custom Pascal scripting for service management.

---

### Task 1: Build Release Binaries
**Files:**
- `bin/prometheus.exe` (Go binary)
- `bin/ghost-shell.exe` (Tauri binary)

- [ ] **Step 1.1: Build Go binary in release mode**
  ```
  go build -ldflags="-s -w" -o bin/prometheus.exe ./cmd/prometheus
  ```
- [ ] **Step 1.2: Build Ghost Shell via Tauri**
  ```
  cd ghost-shell && npm run tauri build
  ```
  Output: `ghost-shell/src-tauri/target/release/ghost-shell.exe`
- [ ] **Step 1.3: Copy Tauri binary to `bin/`**
- [ ] **Step 1.4: Verify both binaries exist and are valid PE files**

### Task 2: Enhanced Inno Setup Script
**Files:**
- Modify: `scripts/windows/prometheus.iss`

- [ ] **Step 2.1: Update [Setup] section**
  - AppVersion = 2.0.0 (major version bump for new architecture)
  - AppPublisher = "Kora AI Lab"
  - DefaultDirName = `{localappdata}\Programs\Prometheus`
  - OutputDir = `release/`
  - OutputBaseFilename = `prometheus-2.0.0-setup`

- [ ] **Step 2.2: Update [Files] section**
  - Include `bin/prometheus.exe` → `{app}\prometheus.exe`
  - Include `bin/ghost-shell.exe` → `{app}\prometheus-shell.exe`
  - Include Ghost Shell assets (www folder if needed)

- [ ] **Step 2.3: Add [Run] section for service installation**
  ```
  [Run]
  Filename: "{app}\prometheus.exe"; Parameters: "service install"; Flags: runhidden waituntilterminated
  Filename: "{app}\prometheus.exe"; Parameters: "service start"; Flags: runhidden waituntilterminated
  Filename: "{app}\prometheus-shell.exe"; Description: "Launch Ghost Shell"; Tasks: launchshell; Flags: nowait postinstall skipifsilent
  ```

- [ ] **Step 2.4: Add [UninstallRun] section for service removal**
  ```
  [UninstallRun]
  Filename: "{app}\prometheus.exe"; Parameters: "service stop"; Flags: runhidden waituntilterminated
  Filename: "{app}\prometheus.exe"; Parameters: "service uninstall"; Flags: runhidden waituntilterminated
  ```

- [ ] **Step 2.5: Update [Icons] section**
  - Start Menu: "Prometheus Ghost Shell" → `prometheus-shell.exe`
  - Desktop: optional shortcut
  - Start Menu: "Prometheus Service Manager" → `prometheus.exe` (for CLI access)
  - Uninstall shortcut

- [ ] **Step 2.6: Add error handling in [Code]**
  - Check if service install succeeds, show warning if not
  - Detect if service is already installed, offer upgrade path

### Task 3: Build Script
**Files:**
- Modify: `scripts/windows/build-msi.ps1`

- [ ] **Step 3.1: Update build script to:**
  1. Build Go binary (`go build -ldflags="-s -w"`)
  2. Build Tauri binary (`npm run tauri build`)
  3. Copy binaries to `release/` folder
  4. Run Inno Setup compiler (`ISCC.exe scripts/windows/prometheus.iss`)
  5. Output: `release/prometheus-2.0.0-setup.exe`

- [ ] **Step 3.2: Add version injection**
  - Read version from `go_version/cmd/prometheus/main.go` or a VERSION file
  - Inject into Inno Setup script dynamically

### Task 4: Verification
- [ ] **Step 4.1: Run build script**
  - Verify MSI/EXE is produced
  - Check file size is reasonable (< 50MB)
- [ ] **Step 4.2: Test installation (manual)**
  - Run installer
  - Verify service is installed (`sc query PrometheusCore`)
  - Verify Ghost Shell launches
  - Verify Start Menu shortcuts
- [ ] **Step 4.3: Test uninstallation (manual)**
  - Run uninstaller
  - Verify service is removed
  - Verify files are cleaned up

---

## Prerequisites
- **Inno Setup** installed (or use `iscc` from PATH)
- **Go toolchain** (already available)
- **Rust/Cargo** (already available)
- **Node.js** (already available)

## Constraints
- Installer must work on Windows 10/11 x64
- Must handle upgrade from v1.x (CLI-only) to v2.x (service + shell)
- Service install requires admin privileges
- Ghost Shell does NOT require admin (runs as user)
