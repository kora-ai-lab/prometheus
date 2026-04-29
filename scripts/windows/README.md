# Windows Installer Build Instructions

## Local Build with Inno Setup

### Prerequisites
1. Download and install Inno Setup from: https://jrsoftware.org/isdl.php
2. Build the Prometheus binary:
   ```powershell
   cd go_version
   go build -trimpath -ldflags="-s -w" -o prometheus.exe ./cmd/prometheus
   ```

### Build Installer
```powershell
cd scripts\windows
# Copy binary to current directory
Copy-Item ..\..\prometheus.exe prometheus-windows-amd64.exe
# Build with Inno Setup
& "C:\Program Files (x86)\Inno Setup 6\ISCC.exe" prometheus.iss
```

The installer will be created in `Output\prometheus-windows-amd64-setup.exe`

## CI Build

The GitHub Actions workflow automatically:
1. Builds the Windows binary
2. Installs Inno Setup via Chocolatey
3. Builds the installer
4. Attaches it to the GitHub release

## Installer Features

- Installs to `C:\Program Files\Prometheus`
- Adds to PATH automatically
- Creates Start Menu shortcuts
- Optional desktop shortcut
- Uninstaller via Control Panel
- Professional Windows metadata

## Testing

After installation:
```powershell
# Test PATH
prometheus --help

# Test web UI
prometheus --web

# Test double-click (should launch web UI)
```

## Troubleshooting

If PATH not updated:
- Close and reopen terminal
- Or manually add: `setx PATH "%PATH%;C:\Program Files\Prometheus"`
