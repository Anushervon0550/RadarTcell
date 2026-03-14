package service

import (
	"errors"
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

func TestValidateTechUpsert_RejectsNonUUIDMetricID(t *testing.T) {
	err := validateTechUpsert(domain.TechnologyUpsert{
		Slug:      "sample-tech",
		Name:      "Sample",
		TrendSlug: "trend-ai",
		TRL:       5,
		Index:     1,
		CustomMetrics: []domain.TechnologyMetricValueUpsert{
			{MetricID: "not-a-uuid"},
		},
	}, true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected domain.ErrInvalid, got %v", err)
	}
}

func TestValidateTechUpsert_AllowsUUIDMetricID(t *testing.T) {
	err := validateTechUpsert(domain.TechnologyUpsert{
		Slug:      "sample-tech",
		Name:      "Sample",
		TrendSlug: "trend-ai",
		TRL:       5,
		Index:     1,
		CustomMetrics: []domain.TechnologyMetricValueUpsert{
			{MetricID: "3f7fd8ec-6e11-4a1a-a4ab-c4db06f94d89"},
		},
	}, true)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

