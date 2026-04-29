# Prometheus Core Service Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transform the Prometheus CLI into a headless Windows service with an asynchronous API for the Ghost Shell.

**Architecture:**
A long-running background process that exposes a REST API on `localhost:8080`. Instead of blocking requests, it will use a Task ID system where the UI can poll or stream status updates.

**Tech Stack:** Go 1.25, `golang.org/x/sys/windows` (for service), `net/http`.

**Current State:**
- `cmd/prometheus/main.go` (260 lines): CLI entry point with `--web` flag that starts a blocking HTTP server
- `internal/ui/webui.go` (119 lines): Basic HTTP server with synchronous `/api/execute`
- `internal/task/task.go` (67 lines): Task struct with ID, Goal, Status, Context, Memory, BlockedReason
- `internal/task/loop.go` (202 lines): Task.Run() loop that executes LLM actions
- `internal/task/deps.go` (28 lines): TaskDeps struct with all dependencies

---

### Task 1: Async Task Manager
**Files:**
- Modify: `internal/task/task.go` - Add `Result`, `Progress`, `Error` fields
- Create: `internal/api/manager.go` - TaskManager with concurrent map
- Create: `internal/api/manager_test.go` - Basic tests

- [ ] **Step 1.1: Extend Task struct** (`internal/task/task.go`)
  - Add `Result string` - final output of completed task
  - Add `Progress string` - current status message ("Initializing...", "Executing command...", "Done")
  - Add `Error string` - error message if failed
  - Add `UpdatedAt time.Time` already exists, ensure it's updated on progress changes
  - Add `SetProgress(msg string)` method

- [ ] **Step 1.2: Create TaskManager** (`internal/api/manager.go`)
  ```go
  type TaskManager struct {
      mu       sync.RWMutex
      tasks    map[string]*task.Task
      active   map[string]context.CancelFunc // to support cancellation
      newDeps  func() *task.TaskDeps         // factory for task dependencies
  }
  ```
  - `NewTaskManager(depsFactory func() *task.TaskDeps) *TaskManager`
  - `Submit(goal string) string` - creates task, starts goroutine, returns ID
  - `GetStatus(id string) (*task.Task, error)` - returns task state
  - `ListActive() []*task.Task` - returns all non-terminal tasks
  - `Cancel(id string) error` - cancels a running task
  - Thread-safe with `sync.RWMutex`

- [ ] **Step 1.3: Wire Task.Run to update progress**
  - Modify `internal/task/loop.go` to accept a progress callback or update `t.Progress` at key points:
    - Before LLM call: "Thinking..."
    - After action parse: "Executing: {action}"
    - On blocked: "Waiting for input: {reason}"
    - On done/failed: terminal state
  - Call `deps.TaskStore.Save(t)` after progress updates

- [ ] **Step 1.4: Commit**
  `git commit -m "feat: implement async task manager with progress tracking"`

### Task 2: Headless API Server Refactor
**Files:**
- Modify: `internal/ui/webui.go` - Complete API overhaul
- Modify: `assets/static/index.html` - Update to use async API

- [ ] **Step 2.1: Refactor WebServer to accept TaskManager**
  - Change constructor: `NewWebServer(host string, port int, mgr *api.TaskManager, cfg *config.Config)`
  - Remove old `interface{}` parameters

- [ ] **Step 2.2: Update `/api/execute` to async**
  - Accept `POST {"goal": "..."}` 
  - Call `mgr.Submit(goal)` 
  - Return `{"task_id": "..."}` immediately (HTTP 202 Accepted)

- [ ] **Step 2.3: Add `/api/tasks/{id}` endpoint**
  - `GET /api/tasks/{id}` - returns task status, progress, result
  - `DELETE /api/tasks/{id}` - cancel task
  - `GET /api/tasks` - list all tasks (with optional `?status=active` filter)

- [ ] **Step 2.4: Add SSE streaming endpoint**
  - `GET /api/tasks/{id}/stream` - Server-Sent Events for real-time progress
  - Sends `data: {"progress": "...", "status": "..."}` on each update
  - Closes when task reaches terminal state

- [ ] **Step 2.5: Add API token authentication**
  - Read token from config or generate on first start
  - Check `Authorization: Bearer <token>` header on all API endpoints
  - Skip auth for `GET /` (static assets)

- [ ] **Step 2.6: Add `/api/health` endpoint**
  - Returns `{"status": "ok", "version": "1.0.3", "uptime": "..."}`

- [ ] **Step 2.7: Commit**
  `git commit -m "feat: refactor webui to async API with SSE streaming"`

### Task 3: Windows Service Integration
**Files:**
- Create: `internal/service/service.go` - Cross-platform service abstraction
- Create: `internal/service/windows.go` - Windows-specific implementation (build tag)
- Create: `internal/service/unix.go` - Unix stub (build tag)
- Modify: `cmd/prometheus/main.go` - Add service subcommands

- [ ] **Step 3.1: Add `golang.org/x/sys/windows` dependency**
  - Run `go get golang.org/x/sys/windows`
  - Also need `golang.org/x/sys/windows/svc` for service management

- [ ] **Step 3.2: Implement service installation** (`internal/service/windows.go`)
  - `Install() error` - registers service with SCM using `mgr.Create()`
  - `Uninstall() error` - removes service
  - `Status() (string, error)` - queries service state
  - Use `golang.org/x/sys/windows/svc/mgr` package
  - Set service to `StartType: automatic`
  - Configure recovery: restart on failure

- [ ] **Step 3.3: Implement service runner**
  - `RunService(ctx context.Context, mgr *api.TaskManager, srv *ui.WebServer) error`
  - Implements `windows/service.Handler` interface
  - Handles `svc.Start`, `svc.Stop`, `svc.Interrogate`
  - Reports status to SCM: `svc.Running`, `svc.StopPending`

- [ ] **Step 3.4: Refactor main.go**
  - Add subcommands: `service install`, `service uninstall`, `service status`, `service start`, `service stop`
  - Add `--service` flag to run as service directly (used by SCM)
  - Keep existing CLI behavior when no service flags
  - Version bump to v1.0.3

- [ ] **Step 3.5: Commit**
  `git commit -m "feat: add windows service management with auto-start"`

### Task 4: Core Service Entrypoint
**Files:**
- Create: `cmd/prometheus-service/main.go` - Dedicated service binary (optional)

- [ ] **Step 4.1: Create dedicated service entrypoint**
  - Alternative to `--service` flag: separate binary for cleaner separation
  - Same dependency initialization as CLI but starts TaskManager + WebServer
  - Logs to file instead of stdout when running as service

- [ ] **Step 4.2: Update build scripts**
  - Modify `Makefile` to build both `prometheus.exe` and `prometheus-service.exe`
  - Update `.github/workflows/release.yml` if needed

- [ ] **Step 4.3: Commit**
  `git commit -m "feat: add dedicated service binary entrypoint"`

### Task 5: Verification & Testing
- [ ] **Step 5.1: Unit tests for TaskManager**
  - Test concurrent task submission
  - Test task cancellation
  - Test progress updates

- [ ] **Step 5.2: Integration test for API**
  - Test full async flow: submit → poll → complete
  - Test SSE streaming
  - Test auth rejection without token

- [ ] **Step 5.3: Manual service test**
  - `prometheus.exe service install`
  - Verify in `services.msc` or `sc query PrometheusCore`
  - `sc start PrometheusCore`
  - `curl http://localhost:8080/api/health`
  - `curl -X POST http://localhost:8080/api/execute -d '{"goal":"test"}'`
  - `sc stop PrometheusCore`
  - `prometheus.exe service uninstall`

- [ ] **Step 5.4: Final Commit**
  `git commit -m "test: add unit and integration tests for core service"`

---

## Dependencies Graph
```
Task 1 (Task Manager) → Task 2 (API Server) → Task 5 (Testing)
                    ↘ Task 3 (Windows Service) → Task 5
Task 4 (Service Binary) → Task 5
```

## Risk Mitigation
- **Thread safety**: Use `sync.RWMutex` for all shared state in TaskManager
- **Context propagation**: Ensure all goroutines respect context cancellation
- **Service recovery**: Configure Windows SCM to auto-restart on crash
- **Backward compat**: Keep `--web` flag working during transition
