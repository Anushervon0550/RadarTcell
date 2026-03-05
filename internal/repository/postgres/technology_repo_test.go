package postgres

import (
	"strings"
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

func TestBuildTechWhere_Filters(t *testing.T) {
	p := domain.TechnologyListParams{
		Search:    "ai",
		TrendID:   "trend-id",
		TagID:     "tag-id",
		HasTRLMin: true,
		TRLMin:    3,
		OnlyIDs:   []string{"1", "2"},
	}

	where, args := buildTechWhere(p, 0)
	if !strings.Contains(where, "tech.name ILIKE") {
		t.Fatalf("expected search filter, got: %s", where)
	}
	if !strings.Contains(where, "tech.trend_id = $2::uuid") {
		t.Fatalf("expected trend filter with $2, got: %s", where)
	}
	if !strings.Contains(where, "technology_tags") {
		t.Fatalf("expected tag filter, got: %s", where)
	}
	if !strings.Contains(where, "readiness_level >= $4") {
		t.Fatalf("expected TRL filter, got: %s", where)
	}
	if !strings.Contains(where, "ANY($5::text[])") {
		t.Fatalf("expected only_ids filter, got: %s", where)
	}

	if len(args) != 5 {
		t.Fatalf("expected 5 args, got %d", len(args))
	}
	if args[0] != "ai" || args[1] != "trend-id" || args[2] != "tag-id" || args[3] != 3 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestNormalizeOrder(t *testing.T) {
	if normalizeOrder("desc") != "DESC" {
		t.Fatal("expected DESC")
	}
	if normalizeOrder("anything") != "ASC" {
		t.Fatal("expected default ASC")
	}
}
