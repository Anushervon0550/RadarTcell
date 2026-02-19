package httpapi

import (
	"net/http"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type CatalogHandler struct {
	svc ports.CatalogService
}

func NewCatalogHandler(svc ports.CatalogService) *CatalogHandler {
	return &CatalogHandler{svc: svc}
}

func (h *CatalogHandler) ListTrends(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListTrends(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CatalogHandler) ListSDGs(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListSDGs(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CatalogHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListTags(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CatalogHandler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListOrganizations(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CatalogHandler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListMetrics(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}
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
