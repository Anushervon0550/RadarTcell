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
	if err := validateMetric(cmd); err != nil {
		return "", err
	}
	return s.repo.Create(ctx, cmd)
}

func (s *AdminMetricService) Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return false, fmt.Errorf("%w: id is required", domain.ErrInvalid)
	}
	if err := validateMetric(cmd); err != nil {
		return false, err
	}
	return s.repo.Update(ctx, id, cmd)
}

func (s *AdminMetricService) Delete(ctx context.Context, id string) (bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return false, fmt.Errorf("%w: id is required", domain.ErrInvalid)
	}
	return s.repo.Delete(ctx, id)
}

func validateMetric(cmd domain.MetricDefinitionUpsert) error {
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	t := strings.TrimSpace(cmd.Type)
	if t != "bubble" && t != "bar" {
		return fmt.Errorf("%w: type must be bubble|bar", domain.ErrInvalid)
	}
	return nil
}
