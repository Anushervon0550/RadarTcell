package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type fakeAdminMetricRepo struct {
	createFn func(ctx context.Context, cmd domain.MetricDefinitionUpsert) (string, error)
	updateFn func(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (bool, error)
	deleteFn func(ctx context.Context, id string) (bool, error)

	lastCreateCmd domain.MetricDefinitionUpsert
	lastUpdateID  string
	lastUpdateCmd domain.MetricDefinitionUpsert
	lastDeleteID  string
}

func (f *fakeAdminMetricRepo) Create(ctx context.Context, cmd domain.MetricDefinitionUpsert) (string, error) {
	f.lastCreateCmd = cmd
	if f.createFn != nil {
		return f.createFn(ctx, cmd)
	}
	return "metric-1", nil
}

func (f *fakeAdminMetricRepo) Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (bool, error) {
	f.lastUpdateID = id
	f.lastUpdateCmd = cmd
	if f.updateFn != nil {
		return f.updateFn(ctx, id, cmd)
	}
	return true, nil
}

func (f *fakeAdminMetricRepo) Delete(ctx context.Context, id string) (bool, error) {
	f.lastDeleteID = id
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return true, nil
}

func strPtr(s string) *string { return &s }

func TestAdminMetricService_Create_AllowsDistanceAndFieldKey(t *testing.T) {
	repo := &fakeAdminMetricRepo{}
	svc := NewAdminMetricService(repo, nil)

	id, err := svc.Create(context.Background(), domain.MetricDefinitionUpsert{
		Name:        "  Technology Readiness Level  ",
		Type:        "  distance ",
		Description: strPtr("TRL"),
		Orderable:   true,
		FieldKey:    strPtr(" readiness_level "),
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id == "" {
		t.Fatal("expected id, got empty")
	}

	// Проверяем, что сервис нормализовал данные перед передачей в repo
	if repo.lastCreateCmd.Name != "Technology Readiness Level" {
		t.Fatalf("expected trimmed name, got %q", repo.lastCreateCmd.Name)
	}
	if repo.lastCreateCmd.Type != "distance" {
		t.Fatalf("expected type=distance, got %q", repo.lastCreateCmd.Type)
	}
	if repo.lastCreateCmd.FieldKey == nil || *repo.lastCreateCmd.FieldKey != "readiness_level" {
		t.Fatalf("expected field_key=readiness_level, got %#v", repo.lastCreateCmd.FieldKey)
	}
}

func TestAdminMetricService_Create_RejectsInvalidType(t *testing.T) {
	repo := &fakeAdminMetricRepo{}
	svc := NewAdminMetricService(repo, nil)

	_, err := svc.Create(context.Background(), domain.MetricDefinitionUpsert{
		Name:      "Bad Metric",
		Type:      "pie",
		Orderable: true,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected domain.ErrInvalid, got %v", err)
	}
}

func TestAdminMetricService_Create_RejectsInvalidFieldKey(t *testing.T) {
	repo := &fakeAdminMetricRepo{}
	svc := NewAdminMetricService(repo, nil)

	_, err := svc.Create(context.Background(), domain.MetricDefinitionUpsert{
		Name:      "Bad FieldKey Metric",
		Type:      "distance",
		Orderable: true,
		FieldKey:  strPtr("not_allowed_field"),
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected domain.ErrInvalid, got %v", err)
	}
}

func TestAdminMetricService_Update_RejectsEmptyID(t *testing.T) {
	repo := &fakeAdminMetricRepo{}
	svc := NewAdminMetricService(repo, nil)

	ok, err := svc.Update(context.Background(), "   ", domain.MetricDefinitionUpsert{
		Name:      "Any",
		Type:      "bar",
		Orderable: true,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if ok {
		t.Fatal("expected ok=false")
	}
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected domain.ErrInvalid, got %v", err)
	}
}

func TestAdminMetricService_Update_NormalizesFields(t *testing.T) {
	repo := &fakeAdminMetricRepo{}
	svc := NewAdminMetricService(repo, nil)

	ok, err := svc.Update(context.Background(), " metric-id ", domain.MetricDefinitionUpsert{
		Name:      "  List Index Metric API  ",
		Type:      " distance ",
		Orderable: true,
		FieldKey:  strPtr(" list_index "),
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}

	if repo.lastUpdateID != "metric-id" {
		t.Fatalf("expected trimmed id, got %q", repo.lastUpdateID)
	}
	if repo.lastUpdateCmd.Name != "List Index Metric API" {
		t.Fatalf("expected trimmed name, got %q", repo.lastUpdateCmd.Name)
	}
	if repo.lastUpdateCmd.Type != "distance" {
		t.Fatalf("expected normalized type=distance, got %q", repo.lastUpdateCmd.Type)
	}
	if repo.lastUpdateCmd.FieldKey == nil || *repo.lastUpdateCmd.FieldKey != "list_index" {
		t.Fatalf("expected field_key=list_index, got %#v", repo.lastUpdateCmd.FieldKey)
	}
}
