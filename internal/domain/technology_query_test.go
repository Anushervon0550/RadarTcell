package domain

import "testing"

func TestNormalizeAndValidateTechnologyListParams_Validates(t *testing.T) {
	p := TechnologyListParams{
		TRLMin:    5,
		TRLMax:    3,
		HasTRLMin: true,
		HasTRLMax: true,
	}
	if err := NormalizeAndValidateTechnologyListParams(&p); err == nil {
		t.Fatal("expected error for TRL range")
	}
}

func TestNormalizeAndValidateTechnologyListParams_Locale(t *testing.T) {
	p := TechnologyListParams{Locale: " RU "}
	if err := NormalizeAndValidateTechnologyListParams(&p); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if p.Locale != "ru" {
		t.Fatalf("expected locale=ru, got %q", p.Locale)
	}
}

func TestNormalizeAndValidateTechnologyListParams_Order(t *testing.T) {
	p := TechnologyListParams{Order: "bad"}
	if err := NormalizeAndValidateTechnologyListParams(&p); err == nil {
		t.Fatal("expected error for order")
	}
}

func TestNormalizeAndValidateTechnologyListParams_CursorRequiresListIndexAsc(t *testing.T) {
	p := TechnologyListParams{Cursor: "10:abc", SortBy: "name", Order: "asc"}
	if err := NormalizeAndValidateTechnologyListParams(&p); err == nil {
		t.Fatal("expected error for cursor with non-list_index sort")
	}

	p = TechnologyListParams{Cursor: "10:abc", SortBy: "list_index", Order: "desc"}
	if err := NormalizeAndValidateTechnologyListParams(&p); err == nil {
		t.Fatal("expected error for cursor with non-asc order")
	}
}

func TestNormalizeAndValidateTechnologyListParams_CursorValid(t *testing.T) {
	p := TechnologyListParams{Cursor: "10:abc", SortBy: "list_index", Order: "asc"}
	if err := NormalizeAndValidateTechnologyListParams(&p); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestNormalizeAndValidateTechnologyListParams_SortByUnsupported(t *testing.T) {
	p := TechnologyListParams{SortBy: "commercial_readiness"}
	if err := NormalizeAndValidateTechnologyListParams(&p); err == nil {
		t.Fatal("expected error for unsupported sort_by")
	}
}

func TestNormalizeAndValidateTechnologyListParams_SortBySupported(t *testing.T) {
	p := TechnologyListParams{SortBy: "trend", Order: "desc"}
	if err := NormalizeAndValidateTechnologyListParams(&p); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

