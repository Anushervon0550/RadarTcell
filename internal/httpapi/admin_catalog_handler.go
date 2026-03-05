package httpapi

import (
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type AdminCatalogHandler struct {
	trends ports.AdminTrendService
	tags   ports.AdminTagService
}

func NewAdminCatalogHandler(trends ports.AdminTrendService, tags ports.AdminTagService) *AdminCatalogHandler {
	return &AdminCatalogHandler{trends: trends, tags: tags}
}

// ---- Trends ----

type trendUpsertReq struct {
	Slug        string  `json:"slug,omitempty"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	Order       int     `json:"order_index"`
}

// @Param body body TrendUpsertRequest true "Trend payload"
// @Success 201 {object} IDSlugResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
func (h *AdminCatalogHandler) CreateTrend(w http.ResponseWriter, r *http.Request) {
	var req trendUpsertReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}

	id, err := h.trends.Create(r.Context(), domain.TrendUpsert{
		Slug:        strings.TrimSpace(req.Slug),
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Order:       req.Order,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "slug": strings.TrimSpace(req.Slug)})
}

// @Param slug path string true "Trend slug"
// @Param body body TrendUpsertRequest true "Trend payload"
// @Success 200 {object} IDSlugResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
func (h *AdminCatalogHandler) UpdateTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var req trendUpsertReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}

	id, ok, err := h.trends.Update(r.Context(), slug, domain.TrendUpsert{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Order:       req.Order,
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

// @Success 200 {array} AdminTrendDTO
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *AdminCatalogHandler) ListTrends(w http.ResponseWriter, r *http.Request) {
	items, err := h.trends.List(r.Context())
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// @Param slug path string true "Trend slug"
// @Success 200 {object} AdminTrendDTO
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *AdminCatalogHandler) GetTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	item, ok, err := h.trends.Get(r.Context(), slug)
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

// ---- Tags ----

type tagUpsertReq struct {
	Slug        string  `json:"slug,omitempty"`
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Description *string `json:"description,omitempty"`
}

// @Param body body TagUpsertRequest true "Tag payload"
// @Success 201 {object} IDSlugResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
func (h *AdminCatalogHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req tagUpsertReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}

	id, err := h.tags.Create(r.Context(), domain.TagUpsert{
		Slug:        strings.TrimSpace(req.Slug),
		Title:       req.Title,
		Category:    req.Category,
		Description: req.Description,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "slug": strings.TrimSpace(req.Slug)})
}

// @Param slug path string true "Tag slug"
// @Param body body TagUpsertRequest true "Tag payload"
// @Success 200 {object} IDSlugResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
func (h *AdminCatalogHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var req tagUpsertReq
	if !decodeJSONOr400(w, r, &req) {
		return
	}
	id, ok, err := h.tags.Update(r.Context(), slug, domain.TagUpsert{
		Title:       req.Title,
		Category:    req.Category,
		Description: req.Description,
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

// @Success 200 {array} TagDTO
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *AdminCatalogHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	items, err := h.tags.List(r.Context())
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// @Param slug path string true "Tag slug"
// @Success 200 {object} TagDTO
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *AdminCatalogHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	item, ok, err := h.tags.Get(r.Context(), slug)
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

// @Param slug path string true "Tag slug"
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *AdminCatalogHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	ok, err := h.tags.Delete(r.Context(), slug)
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

// @Param slug path string true "Trend slug"
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *AdminCatalogHandler) DeleteTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	ok, err := h.trends.Delete(r.Context(), slug)
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
