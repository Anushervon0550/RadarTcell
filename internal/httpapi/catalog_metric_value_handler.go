package httpapi

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *CatalogHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	metricID := chi.URLParam(r, "id")
	techID := strings.TrimSpace(r.URL.Query().Get("technology_id"))

	if metricID == "" {
		writeError(w, http.StatusBadRequest, "metric id is required")
		return
	}
	if techID == "" {
		writeError(w, http.StatusBadRequest, "technology_id is required")
		return
	}

	resp, ok, err := h.svc.GetMetricValue(r.Context(), metricID, techID)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
