# Prometheus

An AI-first agent runtime with security, observability, and self-evolution.

## Quick Start

```bash
# Install
curl -fsSL https://raw.githubusercontent.com/prometheus-dev/prometheus/main/scripts/install.sh | sh

# Run
prometheus "Your goal here"
```

## Installation

### Option 1: curl | sh

```bash
curl -fsSL https://raw.githubusercontent.com/prometheus-dev/prometheus/main/scripts/install.sh | sh
```

### Option 2: Build from source

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

## Features

- LLM code generation
- Command execution
- File operations
- Browser automation
- Self-update (`prometheus --update`)