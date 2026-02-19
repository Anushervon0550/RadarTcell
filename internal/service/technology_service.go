package service

import (
	"context"
	"hash/fnv"
	"math"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type TechnologyService struct {
	repo ports.TechnologyRepository
}

func NewTechnologyService(repo ports.TechnologyRepository) *TechnologyService {
	return &TechnologyService{repo: repo}
}

func (s *TechnologyService) GetBySlug(ctx context.Context, slug string) (*domain.Technology, bool, error) {
	return s.repo.GetTechnologyBySlug(ctx, slug)
}

func (s *TechnologyService) List(ctx context.Context, p domain.TechnologyListParams) (domain.TechnologyListResult, error) {
	// highlight -> преобразуем в ID технологий (union)
	if len(p.Highlight) > 0 {
		ids, err := s.resolveHighlight(ctx, p.Highlight)
		if err != nil {
			return domain.TechnologyListResult{}, err
		}
		// можно сохранить в Search/фильтр через repo-реализацию, позже добавим поле OnlyTechnologyIDs в domain
		_ = ids // на следующем шаге подключим в repo
	}

	rows, total, err := s.repo.ListTechnologies(ctx, p)
	if err != nil {
		return domain.TechnologyListResult{}, err
	}

	trendIDs, err := s.repo.ListTrendIDsOrdered(ctx)
	if err != nil {
		return domain.TechnologyListResult{}, err
	}

	trendPos := map[string]int{}
	for i, id := range trendIDs {
		trendPos[id] = i
	}

	segWidth := 2 * math.Pi
	if len(trendIDs) > 0 {
		segWidth = (2 * math.Pi) / float64(len(trendIDs))
	}

	m1min, m1max := minmax(rows, func(t domain.Technology) *float64 { return t.CustomMetric1 })
	m2min, m2max := minmax(rows, func(t domain.Technology) *float64 { return t.CustomMetric2 })
	m3min, m3max := minmax(rows, func(t domain.Technology) *float64 { return t.CustomMetric3 })
	m4min, m4max := minmax(rows, func(t domain.Technology) *float64 { return t.CustomMetric4 })

	items := make([]domain.TechnologyListItem, 0, len(rows))
	for _, t := range rows {
		radius := float64(t.TRL-1) / 8.0
		if radius < 0 {
			radius = 0
		}
		if radius > 1 {
			radius = 1
		}

		u := hashUnit(t.Slug)
		pos := trendPos[t.TrendID]
		angle := float64(pos)*segWidth + u*segWidth

		items = append(items, domain.TechnologyListItem{
			ID:                t.ID,
			Slug:              t.Slug,
			Index:             t.Index,
			Name:              t.Name,
			DescriptionShort:  t.DescriptionShort,
			TRL:               t.TRL,
			TrendID:           t.TrendID,
			TrendSlug:         t.TrendSlug,
			TrendName:         t.TrendName,
			CustomMetric1:     t.CustomMetric1,
			CustomMetric2:     t.CustomMetric2,
			CustomMetric3:     t.CustomMetric3,
			CustomMetric4:     t.CustomMetric4,
			CustomMetric1Norm: norm(t.CustomMetric1, m1min, m1max),
			CustomMetric2Norm: norm(t.CustomMetric2, m2min, m2max),
			CustomMetric3Norm: norm(t.CustomMetric3, m3min, m3max),
			CustomMetric4Norm: norm(t.CustomMetric4, m4min, m4max),
			Angle:             angle,
			Radius:            radius,
		})
	}

	return domain.TechnologyListResult{
		Page:  p.Page,
		Limit: p.Limit,
		Total: total,
		Items: items,
	}, nil
}

func (s *TechnologyService) resolveHighlight(ctx context.Context, tokens []string) ([]string, error) {
	set := map[string]struct{}{}

	add := func(ids []string) {
		for _, id := range ids {
			set[id] = struct{}{}
		}
	}

	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		kind := "tag"
		val := token
		if i := strings.Index(token, ":"); i > 0 {
			kind = strings.ToLower(strings.TrimSpace(token[:i]))
			val = strings.TrimSpace(token[i+1:])
		}

		switch kind {
		case "tag":
			id, ok, err := s.repo.GetTagIDBySlug(ctx, val)
			if err != nil || !ok {
				continue
			}
			ids, err := s.repo.ListTechnologyIDsByTagID(ctx, id)
			if err != nil {
				return nil, err
			}
			add(ids)
		case "trend":
			id, ok, err := s.repo.GetTrendIDBySlug(ctx, val)
			if err != nil || !ok {
				continue
			}
			ids, err := s.repo.ListTechnologyIDsByTrendID(ctx, id)
			if err != nil {
				return nil, err
			}
			add(ids)
		case "sdg":
			id, ok, err := s.repo.GetSDGIDByCode(ctx, val)
			if err != nil || !ok {
				continue
			}
			ids, err := s.repo.ListTechnologyIDsBySDGID(ctx, id)
			if err != nil {
				return nil, err
			}
			add(ids)
		case "organization":
			id, ok, err := s.repo.GetOrganizationIDBySlug(ctx, val)
			if err != nil || !ok {
				continue
			}
			ids, err := s.repo.ListTechnologyIDsByOrganizationID(ctx, id)
			if err != nil {
				return nil, err
			}
			add(ids)
		}
	}

	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out, nil
}

func hashUnit(s string) float64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return float64(h.Sum32()) / float64(^uint32(0))
}

func minmax(rows []domain.Technology, pick func(domain.Technology) *float64) (float64, float64) {
	first := true
	var mn, mx float64
	for _, r := range rows {
		vp := pick(r)
		if vp == nil {
			continue
		}
		v := *vp
		if first {
			mn, mx = v, v
			first = false
			continue
		}
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	if first {
		return 0, 0
	}
	return mn, mx
}

func norm(vp *float64, mn, mx float64) float64 {
	if vp == nil {
		return 0
	}
	if mx <= mn {
		return 1
	}
	return (*vp - mn) / (mx - mn)
}
