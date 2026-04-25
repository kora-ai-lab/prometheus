# Prometheus

Prometheus is a Go-based agent runtime scaffold shaped from the PRD and scaffold documents in this repository.

## Current state

- Core repo structure is in place.
- The executable entrypoint, task loop, config, storage, logging, vault, capability, security, browser, and vision seams are scaffolded.
- `go build ./...` and `go test ./...` pass locally.
- Real embedded `llama-server` artifacts, full browser automation, full vision runtime, and cloud-provider integrations are still explicit placeholders.

## Commands

```powershell
go build ./...
go test ./...
go run ./cmd/prometheus setup
go run ./cmd/prometheus metrics
go run ./cmd/prometheus logs
go run ./cmd/prometheus vault list
go run ./cmd/prometheus selftest
go run ./cmd/prometheus --web
```

