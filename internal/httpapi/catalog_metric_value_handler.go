package httpapi

import (
	"net/http"
)

func (h *CatalogHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	metricID, ok := pathParamRequired(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "metric id is required")
		return
	}

	techID, ok := queryParamRequired(r, "technology_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "technology_id is required")
		return
	}

	resp, found, err := h.svc.GetMetricValue(r.Context(), metricID, techID)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
