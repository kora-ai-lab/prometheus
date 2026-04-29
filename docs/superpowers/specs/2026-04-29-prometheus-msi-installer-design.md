# Design Spec: Prometheus Professional MSI Installer
Date: 2026-04-29
Status: Proposed

## 1. Goal
Provide a professional, one-click installation experience for Windows users that sets up the Core Service and the Ghost Shell UI without requiring manual configuration.

## 2. Installation Flow
The installer will follow a standard Windows MSI flow:
1. **Destination Selection:** Default to `%LOCALAPPDATA%\Programs\Prometheus`.
2. **Core Service Setup:**
    - Extracts the `prometheus-core.exe`.
    - Registers the executable as a Windows Service (`PrometheusCore`).
    - Sets the service to "Automatic" start.
3. **Shell UI Setup:**
    - Installs the Tauri-based Ghost Shell.
    - Creates a Start Menu shortcut.
    - Registers the `Alt + Space` global hotkey.
4. **Dependency Check:** Ensures required C++ redistributables or system components are present.

## 3. Component Map
- **Core Binary:** `prometheus-core.exe` (Go)
- **UI Binary:** `prometheus-shell.exe` (Tauri/Rust)
- **Configuration:** Initial `config.toml` with default local ports.
- **Database:** Initializes empty `prometheus.db` in the application data folder.

## 4. Uninstallation
- Stop and remove the `PrometheusCore` service.
- Remove all files in the installation directory.
- Remove registry keys associated with the app.

## 5. Success Criteria
- User can install via a single `.msi` file.
- Prometheus Core starts automatically on reboot.
- `Alt + Space` immediately opens the Omnibox after installation.
- All assets are correctly placed in `%LOCALAPPDATA%`.
