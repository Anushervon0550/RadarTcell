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
	Slug         string  `json:"slug,omitempty"`
	Name         string  `json:"name"`
	LogoURL      *string `json:"logo_url,omitempty"`
	Description  *string `json:"description,omitempty"`
	Website      *string `json:"website,omitempty"`
	Headquarters *string `json:"headquarters,omitempty"`
}

// @Param body body OrganizationUpsertRequest true "Organization payload"
// @Success 201 {object} IDSlugResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
func (h *AdminOrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req orgUpsertReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}

	id, err := h.svc.Create(r.Context(), domain.OrganizationUpsert{
		Slug:         strings.TrimSpace(req.Slug),
		Name:         req.Name,
		LogoURL:      req.LogoURL,
		Description:  req.Description,
		Website:      req.Website,
		Headquarters: req.Headquarters,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "slug": strings.TrimSpace(req.Slug)})
}

// @Param slug path string true "Organization slug"
// @Param body body OrganizationUpsertRequest true "Organization payload"
// @Success 200 {object} IDSlugResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
func (h *AdminOrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var req orgUpsertReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}

	id, ok, err := h.svc.Update(r.Context(), slug, domain.OrganizationUpsert{
		Name:         req.Name,
		LogoURL:      req.LogoURL,
		Description:  req.Description,
		Website:      req.Website,
		Headquarters: req.Headquarters,
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

// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
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

// @Success 200 {array} OrganizationDTO
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *AdminOrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// @Param slug path string true "Organization slug"
// @Success 200 {object} OrganizationDTO
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *AdminOrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	item, ok, err := h.svc.Get(r.Context(), slug)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}
