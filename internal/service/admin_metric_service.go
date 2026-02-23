package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminMetricService struct {
	repo ports.AdminMetricRepository
}

func NewAdminMetricService(repo ports.AdminMetricRepository) *AdminMetricService {
	return &AdminMetricService{repo: repo}
}

func (s *AdminMetricService) Create(ctx context.Context, cmd domain.MetricDefinitionUpsert) (string, error) {
	if err := validateMetric(&cmd); err != nil {
		return "", err
	}
	fk, err := normalizeAndValidateMetricFieldKey(cmd.FieldKey)
	if err != nil {
		return "", err
	}
	cmd.FieldKey = fk
	return s.repo.Create(ctx, cmd)
}

func (s *AdminMetricService) Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return false, fmt.Errorf("%w: id is required", domain.ErrInvalid)
	}
	if err := validateMetric(&cmd); err != nil {
		return false, err
	}
	fk, err := normalizeAndValidateMetricFieldKey(cmd.FieldKey)
	if err != nil {
		return false, err
	}
	cmd.FieldKey = fk
	return s.repo.Update(ctx, id, cmd)
}

func (s *AdminMetricService) Delete(ctx context.Context, id string) (bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return false, fmt.Errorf("%w: id is required", domain.ErrInvalid)
	}
	return s.repo.Delete(ctx, id)
}

func validateMetric(cmd *domain.MetricDefinitionUpsert) error {
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}

	t := strings.TrimSpace(cmd.Type)
	switch t {
	case "bubble", "bar", "distance":
		// ok
	default:
		return fmt.Errorf("%w: type must be bubble|bar|distance", domain.ErrInvalid)
	}

	cmd.Type = t
	cmd.Name = strings.TrimSpace(cmd.Name)
	return nil
}
func normalizeAndValidateMetricFieldKey(v *string) (*string, error) {
	if v == nil {
		return nil, nil
	}

	s := strings.TrimSpace(*v)
	if s == "" {
		return nil, nil
	}

	switch s {
	case "readiness_level", "list_index",
		"custom_metric_1", "custom_metric_2", "custom_metric_3", "custom_metric_4":
		return &s, nil
	default:
		return nil, fmt.Errorf("%w: field_key must be one of readiness_level, list_index, custom_metric_1..custom_metric_4", domain.ErrInvalid)
	}
}
