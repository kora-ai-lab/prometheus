# Prometheus v2.0.0 Release Notes

## Overview

Prometheus 2.0.0 is a major release of the AI agent runtime platform, featuring a complete rewrite of the core engine with enhanced capabilities, security, and user experience.

## Key Features

### Core Service
- Go-based agent runtime with task loop and configuration management
- Embedded llamaServer integration for local AI inference
- SQLite-based persistent storage layer
- Structured logging with configurable output levels

### Ghost Shell (Tauri UI)
- Modern desktop interface built with Tauri + React
- Real-time agent status monitoring
- Interactive task management console
- System tray integration with quick actions

### Capability Engine
- Extensible plugin architecture for agent capabilities
- Browser automation support (chromedp integration)
- Vision system for screen analysis and OCR
- File system and shell command execution

### Forge
- Agent orchestration and workflow management
- Task scheduling and dependency resolution
- Multi-agent coordination primitives

### Logs System
- Structured JSON logging with rotation
- Configurable log levels (debug, info, warn, error)
- Log aggregation and filtering APIs
- Integration with observability stack

### Security Suite
- Vault-based secret management
- TLS/mTLS support for secure communication
- Capability-based security model
- Sandboxed execution environments

### Observability
- Prometheus metrics endpoint (`/metrics`)
- OpenTelemetry tracing support
- Health check endpoints
- Performance profiling endpoints

### Web UI
- Browser-accessible control panel
- Real-time dashboard with WebSocket updates
- Task execution history and audit logs
- Configuration editor

### Android APK
- Mobile client for remote agent management
- Push notification support
- Offline-capable task queue
- Secure device enrollment

### Windows Installer
- Inno Setup-based installer with silent mode support
- Automatic dependency checking
- Start Menu and Desktop shortcuts
- Clean uninstall with leftover cleanup

## Installation

### Windows (Recommended)
1. Download `prometheus-2.0.0-setup.exe`
2. Run the installer (administrator privileges may be required)
3. Follow the setup wizard or use silent install:
   ```powershell
   prometheus-2.0.0-setup.exe /VERYSILENT /SUPPRESSMSGBOXES
   ```
4. Prometheus will be installed to `C:\ProgramData\Programs\Prometheus\`

### Manual (Any Platform)
```powershell
go build -ldflags="-s -w" -o prometheus.exe ./cmd/prometheus
./prometheus.exe setup
./prometheus.exe --web
```

## System Requirements

- **OS**: Windows 10/11, Linux (amd64/arm64), macOS 12+
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 500MB free space (plus model storage if using local AI)
- **Dependencies**:
  - Chrome/Chromium (for browser automation)
  - SQLite 3.x
  - (Windows) Visual C++ Redistributable

## Known Limitations

- Local AI inference requires manual llamaServer binary placement
- Browser automation limited to Chrome/Chromium (Firefox support pending)
- Vision system requires additional model downloads
- Android APK is experimental and may have connectivity issues
- Cloud provider integrations are placeholder implementations

## Changelog (v1.x → v2.0.0)

### Breaking Changes
- Complete rewrite from Python to Go
- Configuration format changed from YAML to TOML
- API endpoints restructured under `/api/v2/`
- Vault storage format incompatible with v1.x

### New Features
- Ghost Shell Tauri desktop UI
- Capability plugin system
- Built-in metrics and observability
- Windows native installer
- Android mobile client
- Enhanced security with mTLS

### Improvements
- 10x performance improvement over v1.x Python runtime
- Reduced memory footprint (200MB vs 2GB)
- Faster startup time (<1s vs 10s)
- Better error handling and recovery

### Bug Fixes
- Fixed memory leaks in task loop
- Resolved race conditions in capability execution
- Fixed Windows path handling issues
- Corrected timezone handling in logs

## Checksums

```
prometheus.exe (Windows amd64): 25039360 bytes
prometheus-2.0.0-setup.exe:     11121930 bytes
```

## Support

- GitHub Issues: https://github.com/your-org/prometheus/issues
- Documentation: `/docs/`
- Security: See `SECURITY.md`
