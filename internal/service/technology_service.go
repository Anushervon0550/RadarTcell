package service

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type TechnologyService struct {
	repo    ports.TechnologyRepository
	cache   ports.Cache
	listTTL time.Duration
}

func NewTechnologyService(repo ports.TechnologyRepository, cache ports.Cache, listTTL time.Duration) *TechnologyService {
	return &TechnologyService{repo: repo, cache: cache, listTTL: listTTL}
}

func (s *TechnologyService) GetBySlug(ctx context.Context, slug, locale string) (*domain.Technology, bool, error) {
	return s.repo.GetTechnologyBySlug(ctx, slug, locale)
}

func (s *TechnologyService) List(ctx context.Context, p domain.TechnologyListParams) (domain.TechnologyListResult, error) {
	// highlight => фильтр "только выбранные"
	if err := domain.NormalizeAndValidateTechnologyListParams(&p); err != nil {
		var zero domain.TechnologyListResult
		return zero, err
	}
	if len(p.Highlight) > 0 {
		ids, err := s.resolveHighlight(ctx, p.Highlight)
		if err != nil {
			return domain.TechnologyListResult{}, err
		}
		if len(ids) == 0 {
			return domain.TechnologyListResult{
				Page:  p.Page,
				Limit: p.Limit,
				Total: 0,
				Items: []domain.TechnologyListItem{},
			}, nil
		}
		p.OnlyIDs = ids
	}

	if s.cache != nil && s.listTTL > 0 {
		if cached, ok := s.getCachedList(ctx, p); ok {
			return cached, nil
		}
	}

	rows, total, err := s.repo.ListTechnologies(ctx, p)
	if err != nil {
		return domain.TechnologyListResult{}, err
	}

	trendIDs, err := s.repo.ListTrendIDsOrdered(ctx)
	if err != nil {
		return domain.TechnologyListResult{}, err
	}
	trendPos, segWidth := buildTrendPosAndSegWidth(trendIDs)

	// нормализация метрик (0..1) для bubble 01/02 и bar 03/04
	m1min, m1max := minmax(rows, func(t domain.Technology) *float64 { return t.CustomMetric1 })
	m2min, m2max := minmax(rows, func(t domain.Technology) *float64 { return t.CustomMetric2 })
	m3min, m3max := minmax(rows, func(t domain.Technology) *float64 { return t.CustomMetric3 })
	m4min, m4max := minmax(rows, func(t domain.Technology) *float64 { return t.CustomMetric4 })

	items := make([]domain.TechnologyListItem, 0, len(rows))
	for _, t := range rows {
		radius := computeRadius(t.TRL)
		angle := computeAngle(trendPos, segWidth, t.TrendID, t.Slug)

		items = append(items, domain.TechnologyListItem{
			ID:               t.ID,
			Slug:             t.Slug,
			Index:            t.Index,
			Name:             t.Name,
			DescriptionShort: t.DescriptionShort,
			TRL:              t.TRL,

			TrendID:   t.TrendID,
			TrendSlug: t.TrendSlug,
			TrendName: t.TrendName,

			CustomMetric1: t.CustomMetric1,
			CustomMetric2: t.CustomMetric2,
			CustomMetric3: t.CustomMetric3,
			CustomMetric4: t.CustomMetric4,

			CustomMetric1Norm: norm(t.CustomMetric1, m1min, m1max),
			CustomMetric2Norm: norm(t.CustomMetric2, m2min, m2max),
			CustomMetric3Norm: norm(t.CustomMetric3, m3min, m3max),
			CustomMetric4Norm: norm(t.CustomMetric4, m4min, m4max),

			Angle:  angle,
			Radius: radius,
		})
	}

	res := domain.TechnologyListResult{
		Page:  p.Page,
		Limit: p.Limit,
		Total: total,
		Items: items,
	}

	if s.cache != nil && s.listTTL > 0 {
		s.setCachedList(ctx, p, res)
	}
	return res, nil
}

// highlight: множественные значения, тренды/SDG/теги/организации
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

func (s *TechnologyService) GetCard(ctx context.Context, slug, locale string) (domain.TechnologyCard, bool, error) {
	t, ok, err := s.repo.GetTechnologyBySlug(ctx, slug, locale)
	if err != nil || !ok {
		return domain.TechnologyCard{}, ok, err
	}

	trendIDs, err := s.repo.ListTrendIDsOrdered(ctx)
	if err != nil {
		return domain.TechnologyCard{}, false, err
	}
	trendPos, segWidth := buildTrendPosAndSegWidth(trendIDs)

	radius := computeRadius(t.TRL)
	angle := computeAngle(trendPos, segWidth, t.TrendID, t.Slug)

	tags, err := s.repo.ListTagsByTechnologyID(ctx, t.ID)
	if err != nil {
		return domain.TechnologyCard{}, false, err
	}
	sdgs, err := s.repo.ListSDGsByTechnologyID(ctx, t.ID)
	if err != nil {
		return domain.TechnologyCard{}, false, err
	}
	orgs, err := s.repo.ListOrganizationsByTechnologyID(ctx, t.ID)
	if err != nil {
		return domain.TechnologyCard{}, false, err
	}

	return domain.TechnologyCard{
		ID:               t.ID,
		Slug:             t.Slug,
		Index:            t.Index,
		Name:             t.Name,
		DescriptionShort: t.DescriptionShort,
		DescriptionFull:  t.DescriptionFull,
		TRL:              t.TRL,

		TrendID:   t.TrendID,
		TrendSlug: t.TrendSlug,
		TrendName: t.TrendName,

		CustomMetric1: t.CustomMetric1,
		CustomMetric2: t.CustomMetric2,
		CustomMetric3: t.CustomMetric3,
		CustomMetric4: t.CustomMetric4,

		ImageURL:   t.ImageURL,
		SourceLink: t.SourceLink,

		Angle:  angle,
		Radius: radius,

		Tags:          tags,
		SDGs:          sdgs,
		Organizations: orgs,
	}, true, nil
}

func (s *TechnologyService) ListByTrendSlug(ctx context.Context, slug string, p domain.TechnologyListParams) (domain.TechnologyListResult, bool, error) {
	id, ok, err := s.repo.GetTrendIDBySlug(ctx, slug)
	if err != nil || !ok {
		return domain.TechnologyListResult{}, ok, err
	}
	p.TrendID = id
	res, err := s.List(ctx, p)
	return res, true, err
}

func (s *TechnologyService) ListBySDGCode(ctx context.Context, code string, p domain.TechnologyListParams) (domain.TechnologyListResult, bool, error) {
	id, ok, err := s.repo.GetSDGIDByCode(ctx, code)
	if err != nil || !ok {
		return domain.TechnologyListResult{}, ok, err
	}
	p.SDGID = id
	res, err := s.List(ctx, p)
	return res, true, err
}

func (s *TechnologyService) ListByTagSlug(ctx context.Context, slug string, p domain.TechnologyListParams) (domain.TechnologyListResult, bool, error) {
	id, ok, err := s.repo.GetTagIDBySlug(ctx, slug)
	if err != nil || !ok {
		return domain.TechnologyListResult{}, ok, err
	}
	p.TagID = id
	res, err := s.List(ctx, p)
	return res, true, err
}

func (s *TechnologyService) ListByOrganizationSlug(ctx context.Context, slug string, p domain.TechnologyListParams) (domain.TechnologyListResult, bool, error) {
	id, ok, err := s.repo.GetOrganizationIDBySlug(ctx, slug)
	if err != nil || !ok {
		return domain.TechnologyListResult{}, ok, err
	}
	p.OrganizationID = id
	res, err := s.List(ctx, p)
	return res, true, err
}

func (s *TechnologyService) getCachedList(ctx context.Context, p domain.TechnologyListParams) (domain.TechnologyListResult, bool) {
	key := s.techListCacheKey(ctx, p)
	b, ok, err := s.cache.Get(ctx, key)
	if err != nil || !ok {
		return domain.TechnologyListResult{}, false
	}
	var res domain.TechnologyListResult
	if err := json.Unmarshal(b, &res); err != nil {
		return domain.TechnologyListResult{}, false
	}
	return res, true
}

func (s *TechnologyService) setCachedList(ctx context.Context, p domain.TechnologyListParams, res domain.TechnologyListResult) {
	b, err := json.Marshal(res)
	if err != nil {
		return
	}
	_ = s.cache.Set(ctx, s.techListCacheKey(ctx, p), b, s.listTTL)
}

func (s *TechnologyService) techListCacheKey(ctx context.Context, p domain.TechnologyListParams) string {
	version := cacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return "techs:" + version + ":" + encodeTechListParams(p)
}

func encodeTechListParams(p domain.TechnologyListParams) string {
	ids := append([]string(nil), p.OnlyIDs...)
	sort.Strings(ids)

	highlight := append([]string(nil), p.Highlight...)
	sort.Strings(highlight)

	return strings.Join([]string{
		"page=" + itoa(p.Page),
		"limit=" + itoa(p.Limit),
		"search=" + p.Search,
		"trend_id=" + p.TrendID,
		"sdg_id=" + p.SDGID,
		"tag_id=" + p.TagID,
		"org_id=" + p.OrganizationID,
		"trl_min=" + itoa(p.TRLMin),
		"trl_max=" + itoa(p.TRLMax),
		"sort_by=" + p.SortBy,
		"order=" + p.Order,
		"locale=" + p.Locale,
		"highlight=" + strings.Join(highlight, ","),
		"only_ids=" + strings.Join(ids, ","),
	}, "&")
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
