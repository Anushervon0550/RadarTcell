package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminHandler struct {
	auth ports.AuthService
}

func NewAdminHandler(auth ports.AuthService) *AdminHandler {
	return &AdminHandler{auth: auth}
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	token, ok, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(w, http.StatusUnauthorized, "bad credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"token": token})
}

func (h *AdminHandler) Me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"user": AdminSubject(r),
		"role": "admin",
	})
}
