package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type PreferencesHandler struct {
	svc ports.PreferencesService
}

func NewPreferencesHandler(svc ports.PreferencesService) *PreferencesHandler {
	return &PreferencesHandler{svc: svc}
}

type preferencesSaveReq struct {
	UserID   string          `json:"user_id"`
	Settings json.RawMessage `json:"settings"`
}

func (h *PreferencesHandler) Save(w http.ResponseWriter, r *http.Request) {
	var req preferencesSaveReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	req.UserID = strings.TrimSpace(req.UserID)
	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}
	if len(req.Settings) == 0 {
		writeError(w, http.StatusBadRequest, "settings is required")
		return
	}

	p := domain.Preferences{
		UserID:   req.UserID,
		Settings: req.Settings,
	}

	if err := h.svc.Save(r.Context(), p); err != nil {
		writeDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *PreferencesHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := pathParamRequired(r, "user_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "user id is required")
		return
	}
	p, ok, err := h.svc.Get(r.Context(), userID)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id":  p.UserID,
		"settings": p.Settings,
	})
}
