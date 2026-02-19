package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type PreferencesHandler struct {
	svc ports.PreferencesService
}

func NewPreferencesHandler(svc ports.PreferencesService) *PreferencesHandler {
	return &PreferencesHandler{svc: svc}
}

func (h *PreferencesHandler) Save(w http.ResponseWriter, r *http.Request) {
	var p domain.Preferences
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if p.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}
	if err := h.svc.Save(r.Context(), p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *PreferencesHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	p, ok, err := h.svc.Get(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}
