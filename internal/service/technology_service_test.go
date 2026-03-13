package service

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type fakeTechRepo struct {
	listFn      func(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error)
	boundsFn    func(ctx context.Context) (m1min, m1max, m2min, m2max, m3min, m3max, m4min, m4max float64, err error)
	trendIDsFn  func(ctx context.Context) ([]string, error)
	getBySlugFn func(ctx context.Context, slug, locale string) (*domain.Technology, bool, error)
	cardDataFn  func(ctx context.Context, slug, locale string) (domain.TechnologyCardData, bool, error)
	tagsFn      func(ctx context.Context, techID string) ([]domain.Tag, error)
	sdgsFn      func(ctx context.Context, techID string) ([]domain.SDG, error)
	orgsFn      func(ctx context.Context, techID string) ([]domain.Organization, error)
	orderableFn func(ctx context.Context) (map[string]struct{}, error)
}

var _ ports.TechnologyRepository = (*fakeTechRepo)(nil)

func (f *fakeTechRepo) ListTrendIDsOrdered(ctx context.Context) ([]string, error) {
	if f.trendIDsFn != nil {
		return f.trendIDsFn(ctx)
	}
	return []string{}, nil
}

func (f *fakeTechRepo) ListTechnologies(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error) {
	if f.listFn != nil {
		return f.listFn(ctx, p)
	}
	return []domain.Technology{}, 0, nil
}

func (f *fakeTechRepo) GetMetricRanges(ctx context.Context) (m1min, m1max, m2min, m2max, m3min, m3max, m4min, m4max float64, err error) {
	if f.boundsFn != nil {
		return f.boundsFn(ctx)
	}
	return 0, 0, 0, 0, 0, 0, 0, 0, nil
}

func (f *fakeTechRepo) GetTechnologyBySlug(ctx context.Context, slug, locale string) (*domain.Technology, bool, error) {
	if f.getBySlugFn != nil {
		return f.getBySlugFn(ctx, slug, locale)
	}
	return nil, false, nil
}

func (f *fakeTechRepo) GetTechnologyCardDataBySlug(ctx context.Context, slug, locale string) (domain.TechnologyCardData, bool, error) {
	if f.cardDataFn != nil {
		return f.cardDataFn(ctx, slug, locale)
	}
	if f.getBySlugFn != nil {
		t, ok, err := f.getBySlugFn(ctx, slug, locale)
		if err != nil || !ok {
			return domain.TechnologyCardData{}, ok, err
		}
		data := domain.TechnologyCardData{Technology: *t}
		if f.tagsFn != nil {
			if data.Tags, err = f.tagsFn(ctx, t.ID); err != nil {
				return domain.TechnologyCardData{}, false, err
			}
		}
		if f.sdgsFn != nil {
			if data.SDGs, err = f.sdgsFn(ctx, t.ID); err != nil {
				return domain.TechnologyCardData{}, false, err
			}
		}
		if f.orgsFn != nil {
			if data.Organizations, err = f.orgsFn(ctx, t.ID); err != nil {
				return domain.TechnologyCardData{}, false, err
			}
		}
		return data, true, nil
	}
	return domain.TechnologyCardData{}, false, nil
}

func (f *fakeTechRepo) GetTrendIDBySlug(ctx context.Context, slug string) (string, bool, error) {
	return "", false, nil
}

func (f *fakeTechRepo) GetSDGIDByCode(ctx context.Context, code string) (string, bool, error) {
	return "", false, nil
}

func (f *fakeTechRepo) GetTagIDBySlug(ctx context.Context, slug string) (string, bool, error) {
	return "", false, nil
}

func (f *fakeTechRepo) GetOrganizationIDBySlug(ctx context.Context, slug string) (string, bool, error) {
	return "", false, nil
}

func (f *fakeTechRepo) ListTechnologyIDsByTrendID(ctx context.Context, trendID string) ([]string, error) {
	return []string{}, nil
}

func (f *fakeTechRepo) ListTechnologyIDsBySDGID(ctx context.Context, sdgID string) ([]string, error) {
	return []string{}, nil
}

func (f *fakeTechRepo) ListTechnologyIDsByTagID(ctx context.Context, tagID string) ([]string, error) {
	return []string{}, nil
}

func (f *fakeTechRepo) ListTechnologyIDsByOrganizationID(ctx context.Context, orgID string) ([]string, error) {
	return []string{}, nil
}

func (f *fakeTechRepo) ListTagsByTechnologyID(ctx context.Context, techID string) ([]domain.Tag, error) {
	if f.tagsFn != nil {
		return f.tagsFn(ctx, techID)
	}
	return []domain.Tag{}, nil
}

func (f *fakeTechRepo) ListSDGsByTechnologyID(ctx context.Context, techID string) ([]domain.SDG, error) {
	if f.sdgsFn != nil {
		return f.sdgsFn(ctx, techID)
	}
	return []domain.SDG{}, nil
}

func (f *fakeTechRepo) ListOrganizationsByTechnologyID(ctx context.Context, techID string) ([]domain.Organization, error) {
	if f.orgsFn != nil {
		return f.orgsFn(ctx, techID)
	}
	return []domain.Organization{}, nil
}

func TestTechnologyServiceList_ComputesCoordsAndNorms(t *testing.T) {
	m1a := 1.0
	m1b := 3.0
	m2b := 5.0
	m3 := 2.0

	repo := &fakeTechRepo{
		trendIDsFn: func(ctx context.Context) ([]string, error) {
			return []string{"trend-1", "trend-2"}, nil
		},
		boundsFn: func(ctx context.Context) (float64, float64, float64, float64, float64, float64, float64, float64, error) {
			return 1, 3, 5, 5, 2, 2, 0, 0, nil
		},
		listFn: func(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error) {
			return []domain.Technology{
				{
					ID:               "1",
					Slug:             "a",
					Index:            1,
					Name:             "A",
					TRL:              1,
					TrendID:          "trend-1",
					TrendSlug:        "t1",
					TrendName:        "T1",
					CustomMetric1:    &m1a,
					CustomMetric3:    &m3,
					DescriptionShort: nil,
					DescriptionFull:  nil,
				},
				{
					ID:            "2",
					Slug:          "b",
					Index:         2,
					Name:          "B",
					TRL:           9,
					TrendID:       "trend-2",
					TrendSlug:     "t2",
					TrendName:     "T2",
					CustomMetric1: &m1b,
					CustomMetric2: &m2b,
					CustomMetric3: &m3,
				},
			}, 2, nil
		},
	}

	svc := NewTechnologyService(repo, nil, 0)
	res, err := svc.List(context.Background(), domain.TechnologyListParams{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(res.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(res.Items))
	}

	first := res.Items[0]
	if first.Radius != 0 {
		t.Fatalf("expected radius=0, got %v", first.Radius)
	}
	if first.CustomMetric1Norm != 0 {
		t.Fatalf("expected custom_metric_1_norm=0, got %v", first.CustomMetric1Norm)
	}
	if first.CustomMetric3Norm != 1 {
		t.Fatalf("expected custom_metric_3_norm=1 when all equal, got %v", first.CustomMetric3Norm)
	}
	if first.Angle < 0 {
		t.Fatalf("expected angle >= 0, got %v", first.Angle)
	}

	second := res.Items[1]
	if second.Radius != 1 {
		t.Fatalf("expected radius=1, got %v", second.Radius)
	}
	if second.CustomMetric1Norm != 1 {
		t.Fatalf("expected custom_metric_1_norm=1, got %v", second.CustomMetric1Norm)
	}
	if second.CustomMetric2Norm != 1 {
		t.Fatalf("expected custom_metric_2_norm=1, got %v", second.CustomMetric2Norm)
	}
}

func TestTechnologyServiceList_UsesGlobalBoundsNotPageRows(t *testing.T) {
	m1 := 3.0

	repo := &fakeTechRepo{
		trendIDsFn: func(ctx context.Context) ([]string, error) {
			return []string{"trend-1"}, nil
		},
		boundsFn: func(ctx context.Context) (float64, float64, float64, float64, float64, float64, float64, float64, error) {
			return 0, 10, 0, 0, 0, 0, 0, 0, nil
		},
		listFn: func(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error) {
			return []domain.Technology{{
				ID:            "1",
				Slug:          "a",
				Index:         1,
				Name:          "A",
				TRL:           5,
				TrendID:       "trend-1",
				TrendSlug:     "t1",
				TrendName:     "T1",
				CustomMetric1: &m1,
			}}, 1, nil
		},
	}

	svc := NewTechnologyService(repo, nil, 0)
	res, err := svc.List(context.Background(), domain.TechnologyListParams{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(res.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(res.Items))
	}

	if math.Abs(res.Items[0].CustomMetric1Norm-0.3) > 1e-9 {
		t.Fatalf("expected custom_metric_1_norm=0.3 from global bounds, got %v", res.Items[0].CustomMetric1Norm)
	}
}

func TestTechnologyServiceList_GlobalRangesKeepNormStableAcrossPages(t *testing.T) {
	target := 50.0
	low := 0.0
	high := 100.0

	repo := &fakeTechRepo{
		trendIDsFn: func(ctx context.Context) ([]string, error) {
			return []string{"trend-1"}, nil
		},
		boundsFn: func(ctx context.Context) (float64, float64, float64, float64, float64, float64, float64, float64, error) {
			// Глобальные границы датасета.
			return 0, 100, 0, 0, 0, 0, 0, 0, nil
		},
		listFn: func(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error) {
			base := domain.Technology{TrendID: "trend-1", TrendSlug: "t1", TrendName: "T1", TRL: 5}
			switch p.Page {
			case 1:
				return []domain.Technology{
					{ID: "target", Slug: "target", Index: 1, Name: "Target", CustomMetric1: &target, TrendID: base.TrendID, TrendSlug: base.TrendSlug, TrendName: base.TrendName, TRL: base.TRL},
					{ID: "hi", Slug: "hi", Index: 2, Name: "High", CustomMetric1: &high, TrendID: base.TrendID, TrendSlug: base.TrendSlug, TrendName: base.TrendName, TRL: base.TRL},
				}, 3, nil
			case 2:
				return []domain.Technology{
					{ID: "lo", Slug: "lo", Index: 3, Name: "Low", CustomMetric1: &low, TrendID: base.TrendID, TrendSlug: base.TrendSlug, TrendName: base.TrendName, TRL: base.TRL},
					{ID: "target", Slug: "target", Index: 1, Name: "Target", CustomMetric1: &target, TrendID: base.TrendID, TrendSlug: base.TrendSlug, TrendName: base.TrendName, TRL: base.TRL},
				}, 3, nil
			default:
				return []domain.Technology{}, 3, nil
			}
		},
	}

	svc := NewTechnologyService(repo, nil, 0)

	res1, err := svc.List(context.Background(), domain.TechnologyListParams{Page: 1, Limit: 2})
	if err != nil {
		t.Fatalf("page1 list err: %v", err)
	}
	res2, err := svc.List(context.Background(), domain.TechnologyListParams{Page: 2, Limit: 2})
	if err != nil {
		t.Fatalf("page2 list err: %v", err)
	}

	findNorm := func(items []domain.TechnologyListItem, slug string) (float64, bool) {
		for _, it := range items {
			if it.Slug == slug {
				return it.CustomMetric1Norm, true
			}
		}
		return 0, false
	}

	n1, ok := findNorm(res1.Items, "target")
	if !ok {
		t.Fatal("target not found on page 1")
	}
	n2, ok := findNorm(res2.Items, "target")
	if !ok {
		t.Fatal("target not found on page 2")
	}

	if math.Abs(n1-0.5) > 1e-9 || math.Abs(n2-0.5) > 1e-9 {
		t.Fatalf("expected stable norm 0.5 on both pages, got page1=%v page2=%v", n1, n2)
	}
}

func TestTechnologyServiceGetCard_OK(t *testing.T) {
	m1 := 2.0
	repo := &fakeTechRepo{
		cardDataFn: func(ctx context.Context, slug, locale string) (domain.TechnologyCardData, bool, error) {
			return domain.TechnologyCardData{
				Technology: domain.Technology{
					ID:            "tech-1",
					Slug:          "a",
					Index:         7,
					Name:          "A",
					TRL:           5,
					TrendID:       "trend-2",
					TrendSlug:     "t2",
					TrendName:     "Trend 2",
					CustomMetric1: &m1,
				},
				Tags:          []domain.Tag{{ID: "tag-1", Slug: "tag-a", Title: "Tag A"}},
				SDGs:          []domain.SDG{{ID: "sdg-1", Code: "1", Title: "No poverty"}},
				Organizations: []domain.Organization{{ID: "org-1", Slug: "org-a", Name: "Org A"}},
			}, true, nil
		},
		trendIDsFn: func(ctx context.Context) ([]string, error) {
			return []string{"trend-1", "trend-2"}, nil
		},
	}

	svc := NewTechnologyService(repo, nil, 0)
	card, ok, err := svc.GetCard(context.Background(), "a", "en")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
	if card.ID != "tech-1" {
		t.Fatalf("expected card id tech-1, got %q", card.ID)
	}
	if len(card.Tags) != 1 || len(card.SDGs) != 1 || len(card.Organizations) != 1 {
		t.Fatalf("expected 1 tag/sdg/org, got %d/%d/%d", len(card.Tags), len(card.SDGs), len(card.Organizations))
	}
	if card.Radius <= 0 || card.Radius >= 1 {
		t.Fatalf("expected radius in (0,1), got %v", card.Radius)
	}
}

func TestTechnologyServiceGetCard_ErrorsOnRelatedQuery(t *testing.T) {
	expected := errors.New("sdgs query failed")
	repo := &fakeTechRepo{
		cardDataFn: func(ctx context.Context, slug, locale string) (domain.TechnologyCardData, bool, error) {
			return domain.TechnologyCardData{}, false, expected
		},
		trendIDsFn: func(ctx context.Context) ([]string, error) { return []string{"trend-1"}, nil },
	}

	svc := NewTechnologyService(repo, nil, 0)
	_, ok, err := svc.GetCard(context.Background(), "a", "en")
	if ok {
		t.Fatal("expected ok=false on related query error")
	}
	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

