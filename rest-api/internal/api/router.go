package api

import (
	"log/slog"
	"net/http"
	"time"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/snapshot", h.Snapshot)
	mux.HandleFunc("GET /api/system", h.SystemInfo)
	mux.HandleFunc("GET /api/interfaces", h.Interfaces)
	mux.HandleFunc("GET /api/addresses", h.Addresses)
	mux.HandleFunc("GET /api/routes", h.Routes)
	mux.HandleFunc("GET /api/dhcp", h.DHCPLeases)

	return corsMiddleware(loggingMiddleware(mux))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration", time.Since(start),
			"remote", r.RemoteAddr,
		)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
