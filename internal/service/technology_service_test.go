package service

import (
	"context"
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type fakeTechRepo struct {
	listFn      func(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error)
	trendIDsFn  func(ctx context.Context) ([]string, error)
	getBySlugFn func(ctx context.Context, slug, locale string) (*domain.Technology, bool, error)
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

func (f *fakeTechRepo) GetTechnologyBySlug(ctx context.Context, slug, locale string) (*domain.Technology, bool, error) {
	if f.getBySlugFn != nil {
		return f.getBySlugFn(ctx, slug, locale)
	}
	return nil, false, nil
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
	return []domain.Tag{}, nil
}

func (f *fakeTechRepo) ListSDGsByTechnologyID(ctx context.Context, techID string) ([]domain.SDG, error) {
	return []domain.SDG{}, nil
}

func (f *fakeTechRepo) ListOrganizationsByTechnologyID(ctx context.Context, techID string) ([]domain.Organization, error) {
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
