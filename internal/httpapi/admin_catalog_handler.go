package httpapi

import (
	"encoding/json"
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
	Slug  string `json:"slug,omitempty"`
	Name  string `json:"name"`
	Order int    `json:"order_index"`
}

func (h *AdminCatalogHandler) CreateTrend(w http.ResponseWriter, r *http.Request) {
	var req trendUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	id, err := h.trends.Create(r.Context(), domain.TrendUpsert{
		Slug:  strings.TrimSpace(req.Slug),
		Name:  req.Name,
		Order: req.Order,
	})
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "slug": strings.TrimSpace(req.Slug)})
}

func (h *AdminCatalogHandler) UpdateTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var req trendUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	id, ok, err := h.trends.Update(r.Context(), slug, domain.TrendUpsert{
		Name:  req.Name,
		Order: req.Order,
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

// ---- Tags ----

type tagUpsertReq struct {
	Slug        string  `json:"slug,omitempty"`
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Description *string `json:"description,omitempty"`
}

func (h *AdminCatalogHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req tagUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
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

func (h *AdminCatalogHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var req tagUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
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
