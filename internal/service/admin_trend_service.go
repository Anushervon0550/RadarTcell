package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminTrendService struct {
	repo ports.AdminTrendRepository
}

func NewAdminTrendService(repo ports.AdminTrendRepository) *AdminTrendService {
	return &AdminTrendService{repo: repo}
}

func (s *AdminTrendService) Create(ctx context.Context, cmd domain.TrendUpsert) (string, error) {
	if err := validateTrend(cmd, true); err != nil {
		return "", err
	}
	return s.repo.Create(ctx, cmd)
}

func (s *AdminTrendService) Update(ctx context.Context, slug string, cmd domain.TrendUpsert) (string, bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if err := validateTrend(cmd, false); err != nil {
		return "", false, err
	}
	return s.repo.Update(ctx, slug, cmd)
}

func (s *AdminTrendService) Delete(ctx context.Context, slug string) (bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	return s.repo.Delete(ctx, slug)
}

func validateTrend(cmd domain.TrendUpsert, isCreate bool) error {
	if isCreate && strings.TrimSpace(cmd.Slug) == "" {
		return fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	if cmd.Order < 0 {
		return fmt.Errorf("%w: order_index must be >= 0", domain.ErrInvalid)
	}
	return nil
}
