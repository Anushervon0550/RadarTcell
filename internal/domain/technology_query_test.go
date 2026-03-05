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
