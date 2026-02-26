package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type AdminSDGHandler struct {
	svc ports.AdminSDGService
}

func NewAdminSDGHandler(svc ports.AdminSDGService) *AdminSDGHandler {
	return &AdminSDGHandler{svc: svc}
}

type sdgUpsertReq struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

func (h *AdminSDGHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req sdgUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	id, err := h.svc.Create(r.Context(), domain.SDGUpsert{
		Code:        strings.TrimSpace(req.Code),
		Title:       strings.TrimSpace(req.Title),
		Description: req.Description,
		Icon:        req.Icon,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "code": strings.TrimSpace(req.Code)})
}

func (h *AdminSDGHandler) Update(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	var req sdgUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	ok, err := h.svc.Update(r.Context(), code, domain.SDGUpsert{
		Title:       strings.TrimSpace(req.Title),
		Description: req.Description,
		Icon:        req.Icon,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"code": code})
}

func (h *AdminSDGHandler) Delete(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	ok, err := h.svc.Delete(r.Context(), code)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
