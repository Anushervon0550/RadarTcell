package httpapi

import (
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type CatalogHandler struct {
	svc ports.CatalogService
}

func NewCatalogHandler(svc ports.CatalogService) *CatalogHandler {
	return &CatalogHandler{svc: svc}
}

// @Success 200 {array} TrendDTO
// @Failure 500 {object} ErrorResponse
// @Param locale query string false "Locale" example(ru)
func (h *CatalogHandler) ListTrends(w http.ResponseWriter, r *http.Request) {
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))
	items, err := h.svc.ListTrends(r.Context(), locale)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// @Success 200 {array} SDGDTO
// @Failure 500 {object} ErrorResponse
// @Param locale query string false "Locale" example(ru)
func (h *CatalogHandler) ListSDGs(w http.ResponseWriter, r *http.Request) {
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))
	items, err := h.svc.ListSDGs(r.Context(), locale)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// @Success 200 {array} TagDTO
// @Failure 500 {object} ErrorResponse
// @Param locale query string false "Locale" example(ru)
func (h *CatalogHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))
	items, err := h.svc.ListTags(r.Context(), locale)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// @Success 200 {array} OrganizationDTO
// @Failure 500 {object} ErrorResponse
// @Param locale query string false "Locale" example(ru)
func (h *CatalogHandler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))
	items, err := h.svc.ListOrganizations(r.Context(), locale)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// @Success 200 {array} MetricDTO
// @Failure 500 {object} ErrorResponse
// @Param locale query string false "Locale" example(ru)
func (h *CatalogHandler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))
	items, err := h.svc.ListMetrics(r.Context(), locale)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// @Param slug path string true "Organization slug"
// @Success 200 {object} OrganizationDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *CatalogHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	org, ok, err := h.svc.GetOrganizationBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, org)
}
