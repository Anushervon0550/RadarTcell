package httpapi

import (
	"net/http"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type AdminUsersHandler struct {
	svc ports.AdminUserService
}

func NewAdminUsersHandler(svc ports.AdminUserService) *AdminUsersHandler {
	return &AdminUsersHandler{svc: svc}
}

type adminUserCreateReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AdminUsersHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *AdminUsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req adminUserCreateReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}
	id, err := h.svc.Create(r.Context(), domain.AdminUserCreate{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "username": req.Username})
}

func (h *AdminUsersHandler) Activate(w http.ResponseWriter, r *http.Request) {
	h.setActive(w, r, true)
}

func (h *AdminUsersHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	h.setActive(w, r, false)
}

func (h *AdminUsersHandler) setActive(w http.ResponseWriter, r *http.Request, active bool) {
	username := chi.URLParam(r, "username")
	ok, err := h.svc.SetActive(r.Context(), username, active)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

