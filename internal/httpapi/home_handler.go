package httpapi

import (
	"math"
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type HomeHandler struct {
	catalog ports.CatalogService
	tech    ports.TechnologyService
}

func NewHomeHandler(catalog ports.CatalogService, tech ports.TechnologyService) *HomeHandler {
	return &HomeHandler{catalog: catalog, tech: tech}
}

type HomeResponse struct {
	Page   int              `json:"page"`
	Limit  int              `json:"limit"`
	Total  int              `json:"total"`
	Locale string           `json:"locale,omitempty"`
	Trends []HomeTrendBlock `json:"trends"`
}

type HomeTrendBlock struct {
	Slug  string         `json:"slug"`
	Name  string         `json:"name"`
	Items []HomeTechItem `json:"items"`
}

type HomeTechItem struct {
	ID               string  `json:"id"`
	Slug             string  `json:"slug"`
	Name             string  `json:"name"`
	DescriptionShort *string `json:"description_short,omitempty"`
	TRL              int     `json:"trl"`
	Stage            string  `json:"stage"`
	Completion       int     `json:"completion"`
	Angle            float64 `json:"angle"`
	Radius           float64 `json:"radius"`
}

// List returns one payload for the home page: trends with nested technologies.
func (h *HomeHandler) List(w http.ResponseWriter, r *http.Request) {
	p, ok := parseTechListParamsStrict(w, r)
	if !ok {
		return
	}
	locale := strings.TrimSpace(r.URL.Query().Get("locale"))

	trends, err := h.catalog.ListTrends(r.Context(), locale)
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	list, err := h.tech.List(r.Context(), p)
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	groups := make(map[string]*HomeTrendBlock, len(trends))
	ordered := make([]HomeTrendBlock, 0, len(trends)+1)

	for _, t := range trends {
		ordered = append(ordered, HomeTrendBlock{Slug: t.Slug, Name: t.Name, Items: []HomeTechItem{}})
		groups[t.Slug] = &ordered[len(ordered)-1]
	}

	for _, item := range list.Items {
		slug := item.TrendSlug
		if slug == "" {
			slug = "other"
		}
		grp := groups[slug]
		if grp == nil {
			ordered = append(ordered, HomeTrendBlock{Slug: slug, Name: item.TrendName, Items: []HomeTechItem{}})
			groups[slug] = &ordered[len(ordered)-1]
			grp = groups[slug]
		}

		grp.Items = append(grp.Items, HomeTechItem{
			ID:               item.ID,
			Slug:             item.Slug,
			Name:             item.Name,
			DescriptionShort: item.DescriptionShort,
			TRL:              item.TRL,
			Stage:            stageFromTRL(item.TRL),
			Completion:       completionFromItem(item),
			Angle:            item.Angle,
			Radius:           item.Radius,
		})
	}

	writeJSON(w, http.StatusOK, HomeResponse{
		Page:   list.Page,
		Limit:  list.Limit,
		Total:  list.Total,
		Locale: locale,
		Trends: ordered,
	})
}

func stageFromTRL(trl int) string {
	switch {
	case trl <= 3:
		return "idea"
	case trl <= 6:
		return "prototype"
	default:
		return "product"
	}
}

func completionFromItem(item domain.TechnologyListItem) int {
	v := 0.0
	if mv := domain.MetricValueByFieldKey(item.CustomMetrics, "custom_metric_1"); mv != nil {
		v = *mv * 100
	} else if item.CustomMetric1 != nil {
		v = *item.CustomMetric1 * 100
	} else {
		v = item.CustomMetric1Norm * 100
	}
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	return int(math.Round(v))
}


