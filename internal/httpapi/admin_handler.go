package httpapi

import (
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

// @Param body body AdminLoginRequest true "Login payload"
// @Success 200 {object} AdminLoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	token, ok, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeInternalError(w)
		return
	}
	if !ok {
		writeError(w, http.StatusUnauthorized, "bad credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"token": token})
}

// @Success 200 {object} AdminMeResponse
// @Failure 401 {object} ErrorResponse
func (h *AdminHandler) Me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"user": AdminSubject(r),
		"role": "admin",
	})
}
