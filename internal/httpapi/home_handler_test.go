package httpapi

import (
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

func TestCompletionFromItem_PrefersDynamicMetric(t *testing.T) {
	fk := "custom_metric_1"
	dyn := 0.81
	legacy := 0.2
	item := domain.TechnologyListItem{
		TechnologyViewBase: domain.TechnologyViewBase{CustomMetric1: &legacy},
		CustomMetrics: []domain.TechnologyMetricValue{{
			MetricID: "m1",
			FieldKey: &fk,
			Value:    &dyn,
		}},
	}

	got := completionFromItem(item)
	if got != 81 {
		t.Fatalf("expected completion 81 from dynamic metric, got %d", got)
	}
}

func TestCompletionFromItem_FallbacksToLegacy(t *testing.T) {
	legacy := 0.34
	item := domain.TechnologyListItem{TechnologyViewBase: domain.TechnologyViewBase{CustomMetric1: &legacy}}

	got := completionFromItem(item)
	if got != 34 {
		t.Fatalf("expected completion 34 from legacy metric, got %d", got)
	}
}

