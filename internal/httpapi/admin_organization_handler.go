package httpapi

import (
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type AdminOrganizationHandler struct {
	svc ports.AdminOrganizationService
}

func NewAdminOrganizationHandler(svc ports.AdminOrganizationService) *AdminOrganizationHandler {
	return &AdminOrganizationHandler{svc: svc}
}

type orgUpsertReq struct {
	Slug    string  `json:"slug,omitempty"`
	Name    string  `json:"name"`
	LogoURL *string `json:"logo_url,omitempty"`
}

func (h *AdminOrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req orgUpsertReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	id, err := h.svc.Create(r.Context(), domain.OrganizationUpsert{
		Slug:    strings.TrimSpace(req.Slug),
		Name:    req.Name,
		LogoURL: req.LogoURL,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "slug": strings.TrimSpace(req.Slug)})
}

func (h *AdminOrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var req orgUpsertReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	id, ok, err := h.svc.Update(r.Context(), slug, domain.OrganizationUpsert{
		Name:    req.Name,
		LogoURL: req.LogoURL,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "slug": slug})
}

func (h *AdminOrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	ok, err := h.svc.Delete(r.Context(), slug)
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
