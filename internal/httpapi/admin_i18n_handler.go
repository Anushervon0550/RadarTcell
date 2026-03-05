package httpapi

import (
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type AdminI18nHandler struct {
	svc ports.AdminI18nService
}

func NewAdminI18nHandler(svc ports.AdminI18nService) *AdminI18nHandler {
	return &AdminI18nHandler{svc: svc}
}

type trendI18nReq struct {
	Locale      string  `json:"locale"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type technologyI18nReq struct {
	Locale           string  `json:"locale"`
	Name             string  `json:"name"`
	DescriptionShort *string `json:"description_short,omitempty"`
	DescriptionFull  *string `json:"description_full,omitempty"`
}

type metricI18nReq struct {
	Locale      string  `json:"locale"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// Trends i18n

func (h *AdminI18nHandler) UpsertTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var req trendI18nReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}
	if err := h.svc.UpsertTrend(r.Context(), slug, domain.TrendI18nUpsert{
		Locale:      strings.TrimSpace(req.Locale),
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *AdminI18nHandler) GetTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))

	item, ok, err := h.svc.GetTrend(r.Context(), slug, locale)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *AdminI18nHandler) DeleteTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))

	ok, err := h.svc.DeleteTrend(r.Context(), slug, locale)
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

// Technologies i18n

func (h *AdminI18nHandler) UpsertTechnology(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var req technologyI18nReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}
	if err := h.svc.UpsertTechnology(r.Context(), slug, domain.TechnologyI18nUpsert{
		Locale:           strings.TrimSpace(req.Locale),
		Name:             req.Name,
		DescriptionShort: req.DescriptionShort,
		DescriptionFull:  req.DescriptionFull,
	}); err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *AdminI18nHandler) GetTechnology(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))

	item, ok, err := h.svc.GetTechnology(r.Context(), slug, locale)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *AdminI18nHandler) DeleteTechnology(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))

	ok, err := h.svc.DeleteTechnology(r.Context(), slug, locale)
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

// Metrics i18n

func (h *AdminI18nHandler) UpsertMetric(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req metricI18nReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}
	if err := h.svc.UpsertMetric(r.Context(), id, domain.MetricI18nUpsert{
		Locale:      strings.TrimSpace(req.Locale),
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *AdminI18nHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))

	item, ok, err := h.svc.GetMetric(r.Context(), id, locale)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *AdminI18nHandler) DeleteMetric(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))

	ok, err := h.svc.DeleteMetric(r.Context(), id, locale)
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
