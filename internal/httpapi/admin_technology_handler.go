package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type AdminTechnologyHandler struct {
	svc ports.AdminTechnologyService
}

func NewAdminTechnologyHandler(svc ports.AdminTechnologyService) *AdminTechnologyHandler {
	return &AdminTechnologyHandler{svc: svc}
}

type techUpsertReq struct {
	Slug string `json:"slug,omitempty"`

	Index int    `json:"index"`
	Name  string `json:"name"`
	TRL   int    `json:"trl"`

	TrendSlug string `json:"trend_slug"`

	DescriptionShort *string `json:"description_short,omitempty"`
	DescriptionFull  *string `json:"description_full,omitempty"`

	CustomMetric1 *float64 `json:"custom_metric_1,omitempty"`
	CustomMetric2 *float64 `json:"custom_metric_2,omitempty"`
	CustomMetric3 *float64 `json:"custom_metric_3,omitempty"`
	CustomMetric4 *float64 `json:"custom_metric_4,omitempty"`

	ImageURL   *string `json:"image_url,omitempty"`
	SourceLink *string `json:"source_link,omitempty"`

	TagSlugs          []string `json:"tag_slugs,omitempty"`
	SDGCodes          []string `json:"sdg_codes,omitempty"`
	OrganizationSlugs []string `json:"organization_slugs,omitempty"`
}

func (h *AdminTechnologyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req techUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	cmd := toTechUpsert(req)
	id, err := h.svc.Create(r.Context(), cmd)
	if err != nil {
		writeDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "slug": strings.TrimSpace(req.Slug)})
}

func (h *AdminTechnologyHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var req techUpsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	cmd := toTechUpsert(req)
	id, ok, err := h.svc.Update(r.Context(), slug, cmd)
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

func (h *AdminTechnologyHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func toTechUpsert(req techUpsertReq) domain.TechnologyUpsert {
	return domain.TechnologyUpsert{
		Slug:      strings.TrimSpace(req.Slug),
		Index:     req.Index,
		Name:      req.Name,
		TRL:       req.TRL,
		TrendSlug: req.TrendSlug,

		DescriptionShort: req.DescriptionShort,
		DescriptionFull:  req.DescriptionFull,

		CustomMetric1: req.CustomMetric1,
		CustomMetric2: req.CustomMetric2,
		CustomMetric3: req.CustomMetric3,
		CustomMetric4: req.CustomMetric4,

		ImageURL:   req.ImageURL,
		SourceLink: req.SourceLink,

		TagSlugs:          req.TagSlugs,
		SDGCodes:          req.SDGCodes,
		OrganizationSlugs: req.OrganizationSlugs,
	}
}
