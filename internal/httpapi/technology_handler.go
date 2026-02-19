package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type TechnologyHandler struct {
	svc ports.TechnologyService
}

func NewTechnologyHandler(svc ports.TechnologyService) *TechnologyHandler {
	return &TechnologyHandler{svc: svc}
}

func (h *TechnologyHandler) List(w http.ResponseWriter, r *http.Request) {
	p := parseTechListParams(r)

	res, err := h.svc.List(r.Context(), p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}

func parseIntDefault(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func parseHighlights(values []string) []string {
	var out []string
	for _, v := range values {
		parts := strings.Split(v, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
	}
	return out
}
func (h *TechnologyHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	card, ok, err := h.svc.GetCard(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (h *TechnologyHandler) ListByTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	p := parseTechListParams(r)
	res, ok, err := h.svc.ListByTrendSlug(r.Context(), slug, p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "trend not found")
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}

func (h *TechnologyHandler) ListBySDG(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	p := parseTechListParams(r)
	res, ok, err := h.svc.ListBySDGCode(r.Context(), code, p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "sdg not found")
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}

func (h *TechnologyHandler) ListByTag(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	p := parseTechListParams(r)
	res, ok, err := h.svc.ListByTagSlug(r.Context(), slug, p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "tag not found")
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}

func (h *TechnologyHandler) ListByOrganization(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	p := parseTechListParams(r)
	res, ok, err := h.svc.ListByOrganizationSlug(r.Context(), slug, p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "organization not found")
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}
func parseTechListParams(r *http.Request) domain.TechnologyListParams {
	q := r.URL.Query()

	page := parseIntDefault(q.Get("page"), 1)
	limit := parseIntDefault(q.Get("limit"), 20)
	if limit > 200 {
		limit = 200
	}
	if page < 1 {
		page = 1
	}

	p := domain.TechnologyListParams{
		Search:    strings.TrimSpace(q.Get("search")),
		SortBy:    strings.TrimSpace(q.Get("sort_by")),
		Order:     strings.TrimSpace(q.Get("order")),
		Page:      page,
		Limit:     limit,
		Highlight: parseHighlights(q["highlight"]),
	}

	if v := strings.TrimSpace(q.Get("trl_min")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.TRLMin = n
			p.HasTRLMin = true
		}
	}
	if v := strings.TrimSpace(q.Get("trl_max")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.TRLMax = n
			p.HasTRLMax = true
		}
	}

	return p
}
