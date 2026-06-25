package firewall

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Service interface {
	GetFilterRules(ctx context.Context) ([]FilterRule, error)
	GetNATRules(ctx context.Context) ([]NATRule, error)
	GetAddressList(ctx context.Context) ([]AdressListEntry, error)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetFilterRules(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetFilterRules(r.Context())
	if err != nil {
		slog.Error("GetFilterRules", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetNATRules(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetNATRules(r.Context())
	if err != nil {
		slog.Error("GetNATRules", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetAddressList(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetAddressList(r.Context())
	if err != nil {
		slog.Error("GetAddressList", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
