package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type AdminMetricsHandler struct {
	svc ports.AdminMetricService
}

func NewAdminMetricsHandler(svc ports.AdminMetricService) *AdminMetricsHandler {
	return &AdminMetricsHandler{svc: svc}
}

type metricUpsertReq struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description *string `json:"description,omitempty"`
	Orderable   bool    `json:"orderable"`
}

func (h *AdminMetricsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req metricUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	id, err := h.svc.Create(r.Context(), domain.MetricDefinitionUpsert{
		Name:        req.Name,
		Type:        strings.TrimSpace(req.Type),
		Description: req.Description,
		Orderable:   req.Orderable,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (h *AdminMetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req metricUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	ok, err := h.svc.Update(r.Context(), id, domain.MetricDefinitionUpsert{
		Name:        req.Name,
		Type:        strings.TrimSpace(req.Type),
		Description: req.Description,
		Orderable:   req.Orderable,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id})
}

func (h *AdminMetricsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ok, err := h.svc.Delete(r.Context(), id)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
