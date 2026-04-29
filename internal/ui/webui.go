package ui

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/api"
	"github.com/kora-ai-lab/prometheus/internal/config"
	"github.com/kora-ai-lab/prometheus/internal/task"
)

type WebServer struct {
	server    *http.Server
	mgr       *api.TaskManager
	cfg       *config.Config
	authToken string
	startTime time.Time
}

func NewWebServer(host string, port int, mgr *api.TaskManager, cfg *config.Config) *WebServer {
	ws := &WebServer{
		mgr:       mgr,
		cfg:       cfg,
		authToken: generateToken(),
		startTime: time.Now(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", ws.handleHealth)
	mux.HandleFunc("/api/execute", ws.authMiddleware(ws.handleExecute))
	mux.HandleFunc("/api/tasks", ws.authMiddleware(ws.handleListTasks))
	mux.HandleFunc("/api/tasks/", ws.authMiddleware(ws.handleTask))

	fileServer := http.FileServer(http.Dir("assets/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/", fileServer)

	ws.server = &http.Server{Addr: fmt.Sprintf("%s:%d", host, port), Handler: mux}
	return ws
}

func (w *WebServer) AuthToken() string {
	return w.authToken
}

func generateToken() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return hex.EncodeToString(b[:])
	}
	return hex.EncodeToString(b[:])
}

func (w *WebServer) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && (r.URL.Path == "/" || r.URL.Path == "/index.html" || strings.HasPrefix(r.URL.Path, "/static/")) {
			next(rw, r)
			return
		}
		if r.URL.Path == "/api/health" {
			next(rw, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != w.authToken {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(rw).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		next(rw, r)
	}
}

func (w *WebServer) handleHealth(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"status":  "ok",
		"version": "1.0.3",
		"uptime":  time.Since(w.startTime).String(),
	})
}

func (w *WebServer) handleExecute(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct{ Goal string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Goal == "" {
		http.Error(rw, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	taskID := w.mgr.Submit(req.Goal)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
	json.NewEncoder(rw).Encode(map[string]string{"task_id": taskID})
}

func (w *WebServer) handleListTasks(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	statusFilter := r.URL.Query().Get("status")
	var tasks []*task.Task
	if statusFilter == "active" {
		tasks = w.mgr.ListActive()
	} else {
		tasks = w.listAllTasks()
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(tasksToResponse(tasks))
}

func (w *WebServer) listAllTasks() []*task.Task {
	return w.mgr.ListAll()
}

func (w *WebServer) handleTask(rw http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	parts := strings.SplitN(path, "/", 2)

	if len(parts) == 0 || parts[0] == "" {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	taskID := parts[0]

	if len(parts) == 2 && parts[1] == "stream" {
		w.handleStream(rw, r, taskID)
		return
	}

	switch r.Method {
	case http.MethodGet:
		t, err := w.mgr.GetStatus(taskID)
		if err != nil {
			http.Error(rw, `{"error": "task not found"}`, http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(taskToResponse(t))
	case http.MethodDelete:
		if err := w.mgr.Cancel(taskID); err != nil {
			http.Error(rw, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]string{"status": "cancelled"})
	default:
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (w *WebServer) handleStream(rw http.ResponseWriter, r *http.Request, taskID string) {
	t, err := w.mgr.GetStatus(taskID)
	if err != nil {
		http.Error(rw, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	lastProgress := t.Progress
	lastStatus := t.Status

	sendUpdate := func() {
		data, _ := json.Marshal(map[string]string{
			"progress": lastProgress,
			"status":   string(lastStatus),
		})
		fmt.Fprintf(rw, "data: %s\n\n", data)
		flusher.Flush()
	}

	sendUpdate()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(500 * time.Millisecond):
			t, err := w.mgr.GetStatus(taskID)
			if err != nil {
				return
			}

			if t.Progress != lastProgress || t.Status != lastStatus {
				lastProgress = t.Progress
				lastStatus = t.Status
				sendUpdate()
			}

			if isTerminalStatus(t.Status) {
				sendUpdate()
				return
			}
		}
	}
}

func taskToResponse(t *task.Task) map[string]interface{} {
	return map[string]interface{}{
		"id":         t.ID,
		"goal":       t.Goal,
		"status":     string(t.Status),
		"progress":   t.Progress,
		"result":     t.Result,
		"error":      t.Error,
		"created_at": t.CreatedAt,
		"updated_at": t.UpdatedAt,
	}
}

func tasksToResponse(tasks []*task.Task) []map[string]interface{} {
	result := make([]map[string]interface{}, len(tasks))
	for i, t := range tasks {
		result[i] = taskToResponse(t)
	}
	return result
}

func isTerminalStatus(s task.TaskStatus) bool {
	return s == task.StatusDone || s == task.StatusFailed || s == task.StatusCancelled
}

func (w *WebServer) Start() error {
	return w.server.ListenAndServe()
}

func (w *WebServer) Shutdown(ctx context.Context) error {
	if w == nil || w.server == nil {
		return nil
	}
	return w.server.Shutdown(ctx)
}

func StaticPath() string {
	return filepath.Join("assets", "static")
}
