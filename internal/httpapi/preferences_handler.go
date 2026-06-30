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

// @Summary Save preferences
// @Tags preferences
// @Accept json
// @Produce json
// @Param body body PreferencesSaveRequest true "Preferences payload"
// @Success 200 {object} PreferencesSaveResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/preferences [post]
func (h *PreferencesHandler) Save(w http.ResponseWriter, r *http.Request) {
	var req preferencesSaveReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}

	req.UserID = strings.TrimSpace(req.UserID)
	subject := AdminSubject(r)
	if subject == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	// Ресурс принадлежит аутентифицированному субъекту: нельзя писать чужие настройки.
	if req.UserID != "" && req.UserID != subject {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if len(req.Settings) == 0 {
		writeError(w, http.StatusBadRequest, "settings is required")
		return
	}

	p := domain.Preferences{
		UserID:   subject,
		Settings: req.Settings,
	}

	if err := h.svc.Save(r.Context(), p); err != nil {
		writeDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

// @Summary Get preferences by user id
// @Tags preferences
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} PreferencesGetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/preferences/{user_id} [get]
func (h *PreferencesHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := pathParamRequired(r, "user_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "user id is required")
		return
	}
	// Доступ только к собственным настройкам.
	if subject := AdminSubject(r); subject == "" || userID != subject {
		writeError(w, http.StatusForbidden, "forbidden")
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
