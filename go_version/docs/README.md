# Prometheus

An AI-first agent runtime with security, observability, and self-evolution.

## Quick Start

**macOS/Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.sh | sh
prometheus "Your goal here"
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.ps1 | iex
prometheus "Your goal here"
```

## Installation

### Option 1: curl | sh (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.sh | sh
```

### Option 2: MSI Installer (Windows - Recommended)

Download and run `prometheus-windows-amd64.msi` from [GitHub Releases](https://github.com/kora-ai-lab/prometheus/releases/latest). The installer will:
- Install Prometheus to Program Files
- Add it to your PATH automatically
- Create Start Menu shortcuts

### Option 3: PowerShell (Windows)

```powershell
irm https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.ps1 | iex
```

### Option 4: Download .exe directly (Windows)

Download the latest `prometheus-windows-amd64.exe` from [GitHub Releases](https://github.com/kora-ai-lab/prometheus/releases/latest) and add it to your PATH.

> **Note:** Windows Defender may flag the executable as a false positive. This is expected for unsigned binaries. Use the MSI installer for better compatibility, or add an exclusion if needed.

### Option 5: Build from source

```bash
git clone https://github.com/kora-ai-lab/prometheus
cd prometheus/go_version
go build -ldflags="-s -w" -o prometheus ./cmd/prometheus
```

## Configuration

Set API key:

```bash
export GROQ_API_KEY=sk-...
```

Or use the Web UI at http://localhost:8080

## Usage

### CLI

```bash
prometheus "Create a hello world app"
prometheus --web  # Start web UI
```

### Web UI

1. Run `prometheus --web`
2. Open http://localhost:8080
3. Enter your goal

## Security

- Sandboxed execution
- Auto-confirmation for risky commands
- Secrets redacted in logs

## Troubleshooting

### Windows Defender False Positive

If Windows Defender flags Prometheus as malware:

1. **Use the MSI installer** - MSI files have better reputation than raw executables
2. **Add an exclusion:**
   - Open Windows Security → Virus & threat protection → Manage settings
   - Scroll to Exclusions → Add or remove exclusions
   - Add the installation folder (e.g., `C:\Program Files\Prometheus`)
3. **Submit to Microsoft:** Report the false positive at https://www.microsoft.com/en-us/wdsi/filesubmission

### Installation Issues

**"prometheus command not found"**
- Close and reopen your terminal after installation
- Or manually add to PATH: `setx PATH "%PATH%;C:\Program Files\Prometheus"`

**Double-clicking .exe closes immediately**
- Use the MSI installer instead
- Or run from PowerShell: `.\prometheus-windows-amd64.exe --help`

**Web UI not accessible**
- Check if port 8080 is available: `netstat -ano | findstr :8080`
- Use a different port: `prometheus --web --port 9090`

## Features

- LLM code generation
- Command execution
- File operations
- Browser automation
- Self-update (`prometheus --update`)