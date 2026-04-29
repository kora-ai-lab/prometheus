# Design Spec: Prometheus Core Service (Headless Backend)
Date: 2026-04-29
Status: Proposed

## 1. Goal
Transform the Prometheus Go binary from a CLI tool into a persistent background service (headless) that acts as the "brain" for the Ghost Shell UI.

## 2. Architecture
The Core will operate as a standalone process implementing a Local API (REST/JSON) to communicate with the frontend.

### Components
- **Service Wrapper:** Handles Windows Service registration and lifecycle (start, stop, restart).
- **Local API Server:** 
    - `POST /api/execute`: Receives goals from the Omnibox, triggers the agent, and streams status.
    - `GET /api/status`: Returns the current state of the agent (Idle, Working, Blocked).
    - `GET /api/metrics`: Returns system and performance data.
    - `POST /api/settings`: Updates configuration.
- **Task Manager:** Manages the execution queue and persists task history to `prometheus.db`.

## 3. Data Flow
`Ghost Shell (UI)` $\rightarrow$ `HTTP POST /api/execute` $\rightarrow$ `Core Service` $\rightarrow$ `AI Agent Runtime` $\rightarrow$ `OS/Web/Files` $\rightarrow$ `UI Result`

## 4. Windows Integration
- **Installation Path:** `%LOCALAPPDATA%\Programs\Prometheus`
- **Service Name:** `PrometheusCore`
- **Auto-start:** Configured to start on system boot.
- **Heartbeat:** The Shell UI will ping `/api/status` to ensure the backend is alive.

## 5. Success Criteria
- The binary can be started/stopped via Windows Services.
- The UI can trigger a goal and receive a status update without needing a terminal window.
- Zero terminal windows appear during normal operation.
