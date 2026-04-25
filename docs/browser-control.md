# Browser Control

The browser layer is scaffolded around a default CDP path and an optional future Playwright path.

Current status:

- `internal/browser/client.go` defines the contract
- `internal/browser/manager.go` wires browser actions into the task loop
- the concrete CDP implementation is still a placeholder client

