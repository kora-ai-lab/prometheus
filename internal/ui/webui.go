package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

type WebServer struct {
	server        *http.Server
	taskExecutor  interface{}
	metrics       interface{}
	config        interface{}
}

func NewWebServer(host string, port int, executor interface{}, metrics interface{}, config interface{}) *WebServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/execute", handleExecute(executor))
	mux.HandleFunc("/api/metrics", handleMetrics(metrics))
	mux.HandleFunc("/api/settings", handleSettings(config))

	fileServer := http.FileServer(http.Dir("assets/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/", fileServer)

	return &WebServer{
		server:  &http.Server{Addr: fmt.Sprintf("%s:%d", host, port), Handler: mux},
		metrics: metrics,
		config:  config,
	}
}

func handleExecute(executor interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct{ Goal string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		exec, ok := executor.(func(ctx context.Context, goal string) (string, error))
		if !ok {
			http.Error(w, "executor not configured", http.StatusInternalServerError)
			return
		}

		result, err := exec(context.Background(), req.Goal)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"result": result})
	}
}

func handleMetrics(metrics interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if m, ok := metrics.(interface{ Snapshot() map[string]interface{} }); ok {
			json.NewEncoder(w).Encode(m.Snapshot())
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"error": "no metrics"})
	}
}

func handleSettings(config interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			if c, ok := config.(interface{ Get(key string) string }); ok {
				json.NewEncoder(w).Encode(map[string]string{
					"model": c.Get("model"),
				})
				return
			}
		case http.MethodPost:
			var req struct{ Key, Value string }
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if c, ok := config.(interface{ Set(key, value string) }); ok {
				c.Set(req.Key, req.Value)
			}
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
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