package httpapi

import "net/http"

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// Healthz godoc
// @Summary Health check
// @Description Returns service liveness status
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthz [get]
func (h *SystemHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}
