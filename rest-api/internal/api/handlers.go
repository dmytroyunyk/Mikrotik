package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/dmytroyunyk/MikrotikApi/internal/mikrotik"
)

type Handler struct {
	svc *mikrotik.Service
	mt  *mikrotik.Client
}

func NewHandler(svc *mikrotik.Service, mt *mikrotik.Client) *Handler {
	return &Handler{svc: svc, mt: mt}
}

type apiResponse struct {
	Data      any       `json:"data,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiResponse{
		Data:      data,
		Timestamp: time.Now(),
	})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiResponse{
		Error:     msg,
		Timestamp: time.Now(),
	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	code := http.StatusOK

	if err := h.mt.Ping(r.Context()); err != nil {
		status = "unhealthy: " + err.Error()
		code = http.StatusServiceUnavailable
	}

	writeJSON(w, code, map[string]string{"status": status})
}

func (h *Handler) Snapshot(w http.ResponseWriter, r *http.Request) {
	snap := h.svc.CollectSnapshot(r.Context())
	writeJSON(w, http.StatusOK, snap)
}

func (h *Handler) SystemInfo(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetSystemResource(r.Context())
	if err != nil {
		slog.Error("SystemInfo", "err", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) Interfaces(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetInterfaces(r.Context())
	if err != nil {
		slog.Error("Interfaces", "err", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) Addresses(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetIPAddresses(r.Context())
	if err != nil {
		slog.Error("Addresses", "err", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) Routes(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetRoutes(r.Context())
	if err != nil {
		slog.Error("Routes", "err", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) DHCPLeases(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetDHCPLeases(r.Context())
	if err != nil {
		slog.Error("DHCPLeases", "err", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}
