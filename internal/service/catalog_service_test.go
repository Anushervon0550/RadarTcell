package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type catalogRepoStub struct {
	ports.CatalogRepository
	getMetricValueFn func(ctx context.Context, metricID, technologyID string) (map[string]any, bool, error)
	gotMetricID      string
	gotTechnologyID  string
}

func (s *catalogRepoStub) GetMetricValue(ctx context.Context, metricID, technologyID string) (map[string]any, bool, error) {
	s.gotMetricID = metricID
	s.gotTechnologyID = technologyID
	if s.getMetricValueFn != nil {
		return s.getMetricValueFn(ctx, metricID, technologyID)
	}
	return nil, false, nil
}

func TestCatalogServiceGetMetricValue_InvalidUUID(t *testing.T) {
	svc := NewCatalogService(&catalogRepoStub{}, nil, 0)

	_, ok, err := svc.GetMetricValue(context.Background(), "bad-id", "also-bad")
	if ok {
		t.Fatal("expected ok=false")
	}
	if err == nil || !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected domain.ErrInvalid, got %v", err)
	}
}

func TestCatalogServiceGetMetricValue_Passthrough(t *testing.T) {
	metricID := "550e8400-e29b-41d4-a716-446655440000"
	techID := "5ab5c441-465e-423b-a940-b47fbfaf088b"
	repo := &catalogRepoStub{
		getMetricValueFn: func(ctx context.Context, m, t string) (map[string]any, bool, error) {
			return map[string]any{"metric_id": m, "technology_id": t}, true, nil
		},
	}
	svc := NewCatalogService(repo, nil, 0)

	res, ok, err := svc.GetMetricValue(context.Background(), metricID, techID)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
	if repo.gotMetricID != metricID || repo.gotTechnologyID != techID {
		t.Fatalf("expected passthrough ids, got metric=%q tech=%q", repo.gotMetricID, repo.gotTechnologyID)
	}
	if res["metric_id"] != metricID {
		t.Fatalf("expected response metric_id=%s, got %v", metricID, res["metric_id"])
	}
}

