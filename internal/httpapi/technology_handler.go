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

// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Param search query string false "Search text"
// @Param sort_by query string false "Sort field"
// @Param order query string false "Sort order (asc|desc)"
// @Param trl_min query int false "Min TRL"
// @Param trl_max query int false "Max TRL"
// @Param trend_id query string false "Trend ID (uuid)"
// @Param sdg_id query string false "SDG ID (uuid/int, depends on impl)"
// @Param tag_id query string false "Tag ID (uuid)"
// @Param organization_id query string false "Organization ID (uuid)"
// @Param highlight query []string false "Highlights (repeatable): tag:ml, trend:ai, organization:openai" collectionFormat(multi)
// @Param locale query string false "Locale" example(ru)
// @Success 200 {object} TechnologyListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *TechnologyHandler) List(w http.ResponseWriter, r *http.Request) {
	p, ok := parseTechListParamsStrict(w, r)
	if !ok {
		return // error already written
	}

	resp, err := h.svc.List(r.Context(), p)
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	// если у resp есть Total (обычно есть) — ставим header как в остальных list endpoints
	// если у тебя поле называется иначе — просто удали эту строку
	w.Header().Set("X-Total-Count", strconv.Itoa(resp.Total))

	writeJSON(w, http.StatusOK, resp)
}

// @Param slug path string true "Technology slug"
// @Success 200 {object} TechnologyDetailDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *TechnologyHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))

	card, ok, err := h.svc.GetCard(r.Context(), slug, locale)
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

// @Param slug path string true "Trend slug"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} TechnologyListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *TechnologyHandler) ListByTrend(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	p, ok := parseTechListParamsStrict(w, r)
	if !ok {
		return
	}

	res, ok2, err := h.svc.ListByTrendSlug(r.Context(), slug, p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok2 {
		writeError(w, http.StatusNotFound, "trend not found")
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}

// @Param code path string true "SDG code (e.g. SDG 09)"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} TechnologyListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *TechnologyHandler) ListBySDG(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	p, ok := parseTechListParamsStrict(w, r)
	if !ok {
		return
	}

	res, ok2, err := h.svc.ListBySDGCode(r.Context(), code, p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok2 {
		writeError(w, http.StatusNotFound, "sdg not found")
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}

// @Param slug path string true "Tag slug"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} TechnologyListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *TechnologyHandler) ListByTag(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	p, ok := parseTechListParamsStrict(w, r)
	if !ok {
		return
	}

	res, ok2, err := h.svc.ListByTagSlug(r.Context(), slug, p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok2 {
		writeError(w, http.StatusNotFound, "tag not found")
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}

// @Param slug path string true "Organization slug"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} TechnologyListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *TechnologyHandler) ListByOrganization(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	p, ok := parseTechListParamsStrict(w, r)
	if !ok {
		return
	}

	res, ok2, err := h.svc.ListByOrganizationSlug(r.Context(), slug, p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok2 {
		writeError(w, http.StatusNotFound, "organization not found")
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(res.Total))
	writeJSON(w, http.StatusOK, res)
}

// -------------------------------
// Strict query parsing + validation
// -------------------------------

func parseTechListParamsStrict(w http.ResponseWriter, r *http.Request) (domain.TechnologyListParams, bool) {
	q := r.URL.Query()

	// defaults
	page := 1
	limit := 20

	// page strict
	if s := strings.TrimSpace(q.Get("page")); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 {
			writeError(w, http.StatusBadRequest, "invalid: page must be >= 1")
			return domain.TechnologyListParams{}, false
		}
		page = n
	}

	// limit strict: 1..200
	if s := strings.TrimSpace(q.Get("limit")); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 200 {
			writeError(w, http.StatusBadRequest, "invalid: limit must be between 1 and 200")
			return domain.TechnologyListParams{}, false
		}
		limit = n
	}

	p := domain.TechnologyListParams{
		Search: strings.TrimSpace(q.Get("search")),
		SortBy: strings.TrimSpace(q.Get("sort_by")),
		Order:  strings.TrimSpace(q.Get("order")),

		Page:  page,
		Limit: limit,

		Highlight: parseHighlights(q["highlight"]),

		TrendID:        strings.TrimSpace(q.Get("trend_id")),
		SDGID:          strings.TrimSpace(q.Get("sdg_id")),
		TagID:          strings.TrimSpace(q.Get("tag_id")),
		OrganizationID: strings.TrimSpace(q.Get("organization_id")),
		Locale:         strings.TrimSpace(q.Get("locale")),
	}

	// TRL strict using flags
	if s := strings.TrimSpace(q.Get("trl_min")); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 9 {
			writeError(w, http.StatusBadRequest, "invalid: trl_min must be 1..9")
			return domain.TechnologyListParams{}, false
		}
		p.TRLMin = n
		p.HasTRLMin = true
	}

	if s := strings.TrimSpace(q.Get("trl_max")); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 9 {
			writeError(w, http.StatusBadRequest, "invalid: trl_max must be 1..9")
			return domain.TechnologyListParams{}, false
		}
		p.TRLMax = n
		p.HasTRLMax = true
	}

	if p.HasTRLMin && p.HasTRLMax && p.TRLMin > p.TRLMax {
		writeError(w, http.StatusBadRequest, "invalid: trl_min must be <= trl_max")
		return domain.TechnologyListParams{}, false
	}

	// sort_by/order strict (если передали — проверяем; если пусто — сервис может поставить дефолт)
	if p.Order != "" && !isAllowedOrder(p.Order) {
		writeError(w, http.StatusBadRequest, "invalid: order must be asc|desc")
		return domain.TechnologyListParams{}, false
	}

	return p, true
}

func isAllowedOrder(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "asc", "desc":
		return true
	default:
		return false
	}
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
