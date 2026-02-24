package httpapi

import (
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
	FieldKey    *string `json:"field_key,omitempty"`
}

// @Param body body MetricUpsertRequest true "Metric payload"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
func (h *AdminMetricsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req metricUpsertReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}

	id, err := h.svc.Create(r.Context(), domain.MetricDefinitionUpsert{
		Name:        req.Name,
		Type:        strings.TrimSpace(req.Type),
		Description: req.Description,
		Orderable:   req.Orderable,
		FieldKey:    req.FieldKey,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

// @Param id path string true "Metric ID (UUID)"
// @Param body body MetricUpsertRequest true "Metric payload"
// @Success 200 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
func (h *AdminMetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req metricUpsertReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}

	ok, err := h.svc.Update(r.Context(), id, domain.MetricDefinitionUpsert{
		Name:        req.Name,
		Type:        strings.TrimSpace(req.Type),
		Description: req.Description,
		Orderable:   req.Orderable,
		FieldKey:    req.FieldKey,
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

// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
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
