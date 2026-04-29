package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/api"
	"github.com/kora-ai-lab/prometheus/internal/config"
	"github.com/kora-ai-lab/prometheus/internal/task"
)

func newTestServer() (*WebServer, *api.TaskManager) {
	cfg := &config.Config{}
	mgr := api.NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})
	mgr.WithRunFn(func(ctx context.Context, t *task.Task, deps *task.TaskDeps) error { return nil })
	ws := NewWebServer("localhost", 0, mgr, cfg)
	return ws, mgr
}

func TestHealthEndpoint(t *testing.T) {
	ws, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	ws.handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", resp["status"])
	}
	if resp["version"] != "1.0.3" {
		t.Errorf("expected version '1.0.3', got %q", resp["version"])
	}
}

func TestHealthNoAuth(t *testing.T) {
	ws, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	ws.authMiddleware(ws.handleHealth)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 without auth, got %d", rr.Code)
	}
}

func TestExecuteEndpoint(t *testing.T) {
	ws, mgr := newTestServer()

	mgr.WithRunFn(func(ctx context.Context, t *task.Task, deps *task.TaskDeps) error {
		t.Status = task.StatusDone
		t.Result = "test result"
		return nil
	})

	body := []byte(`{"goal": "test goal"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/execute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ws.AuthToken())
	rr := httptest.NewRecorder()

	ws.handleExecute(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["task_id"] == "" {
		t.Error("expected task_id in response")
	}
}

func TestExecuteNoAuth(t *testing.T) {
	ws, _ := newTestServer()

	body := []byte(`{"goal": "test goal"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/execute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	ws.authMiddleware(ws.handleExecute)(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestGetTask(t *testing.T) {
	ws, mgr := newTestServer()

	taskID := mgr.Submit("test goal")
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/"+taskID, nil)
	req.Header.Set("Authorization", "Bearer "+ws.AuthToken())
	rr := httptest.NewRecorder()

	ws.handleTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["id"] != taskID {
		t.Errorf("expected id %q, got %v", taskID, resp["id"])
	}
}

func TestGetTaskNotFound(t *testing.T) {
	ws, _ := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/nonexistent", nil)
	req.Header.Set("Authorization", "Bearer "+ws.AuthToken())
	rr := httptest.NewRecorder()

	ws.handleTask(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestListTasks(t *testing.T) {
	ws, mgr := newTestServer()

	mgr.Submit("goal 1")
	mgr.Submit("goal 2")
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+ws.AuthToken())
	rr := httptest.NewRecorder()

	ws.handleListTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var tasks []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestAuthToken(t *testing.T) {
	ws, _ := newTestServer()

	token := ws.AuthToken()
	if token == "" {
		t.Error("expected non-empty auth token")
	}

	if len(token) < 32 {
		t.Error("expected token to be reasonably long")
	}
}

func TestStreamEndpoint(t *testing.T) {
	ws, mgr := newTestServer()

	taskID := mgr.Submit("stream test")

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/"+taskID+"/stream", nil)
	req.Header.Set("Authorization", "Bearer "+ws.AuthToken())
	rr := httptest.NewRecorder()

	go func() {
		time.Sleep(100 * time.Millisecond)
		t, _ := mgr.GetStatus(taskID)
		t.Status = task.StatusDone
		t.Progress = "Completed"
	}()

	done := make(chan bool)
	go func() {
		ws.handleTask(rr, req)
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("stream timed out")
	}
}
