package ui

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
)

type WebServer struct {
	server *http.Server
}

func NewWebServer(host string, port int) *WebServer {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("assets/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/", fileServer)
	return &WebServer{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: mux,
		},
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
