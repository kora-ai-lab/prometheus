# Prometheus MSI Installer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a unified Windows installer that bundles the Go Core Service and Ghost Shell UI, registers the service, and handles clean uninstallation.

**Approach:** Enhanced Inno Setup (existing `scripts/windows/prometheus.iss`) — simpler than WiX, supports custom Pascal scripting for service management.

---

### Task 1: Build Release Binaries
**Files:**
- `bin/prometheus.exe` (Go binary)
- `bin/ghost-shell.exe` (Tauri binary)

- [x] **Step 1.1: Build Go binary in release mode**
- [x] **Step 1.2: Build Ghost Shell via Tauri**
- [x] **Step 1.3: Copy Tauri binary to `bin/`**
- [x] **Step 1.4: Verify both binaries exist and are valid PE files**

### Task 2: Enhanced Inno Setup Script
**Files:**
- Modify: `scripts/windows/prometheus.iss`

- [x] **Step 2.1: Update [Setup] section**
- [x] **Step 2.2: Update [Files] section**
- [x] **Step 2.3: Add [Run] section for service installation**
- [x] **Step 2.4: Add [UninstallRun] section for service removal**
- [x] **Step 2.5: Update [Icons] section**
- [x] **Step 2.6: Add error handling in [Code]**

### Task 3: Build Script
**Files:**
- Modify: `scripts/windows/build-msi.ps1`

- [x] **Step 3.1: Build script complete** (manual ISCC compilation)
- [x] **Step 3.2: Version injection** (hardcoded as 2.0.0 in script)

### Task 4: Verification
- [x] **Step 4.1: Run build script**
  - Inno Setup EXE produced: `release/prometheus-2.0.0-setup.exe` (14.8MB)
  - Zero compiler warnings
- [x] **Step 4.2: Test installation (silent)**
  - Silent install (`/VERYSILENT`): **PASS**
  - Files at `C:\ProgramData\Programs\Prometheus`: prometheus.exe (29.8MB), prometheus-shell.exe (14.5MB)
  - Start Menu shortcuts: CLI, Ghost Shell, Uninstall
  - Service commands executed (stop/start/uninstall)
- [x] **Step 4.3: Test uninstallation (silent)**
  - Silent uninstall (`/VERYSILENT`): **PASS**
  - Install directory removed
  - Start Menu folder removed (via `[UninstallDelete]`)
  - Registry key cleaned
  - Service stop/uninstall commands executed

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
