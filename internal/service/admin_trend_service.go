package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminTrendService struct {
	repo  ports.AdminTrendRepository
	cache ports.Cache
}

func NewAdminTrendService(repo ports.AdminTrendRepository, cache ports.Cache) *AdminTrendService {
	return &AdminTrendService{repo: repo, cache: cache}
}

func (s *AdminTrendService) Create(ctx context.Context, cmd domain.TrendUpsert) (string, error) {
	if err := validateTrend(cmd, true); err != nil {
		return "", err
	}
	id, err := s.repo.Create(ctx, cmd)
	if err != nil {
		return "", err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return id, nil
}

func (s *AdminTrendService) Update(ctx context.Context, slug string, cmd domain.TrendUpsert) (string, bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if err := validateTrend(cmd, false); err != nil {
		return "", false, err
	}
	id, ok, err := s.repo.Update(ctx, slug, cmd)
	if err != nil || !ok {
		return id, ok, err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return id, ok, nil
}

func (s *AdminTrendService) Delete(ctx context.Context, slug string) (bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	ok, err := s.repo.Delete(ctx, slug)
	if err != nil || !ok {
		return ok, err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return ok, nil
}

func (s *AdminTrendService) List(ctx context.Context) ([]domain.AdminTrend, error) {
	return s.repo.List(ctx)
}

func (s *AdminTrendService) Get(ctx context.Context, slug string) (domain.AdminTrend, bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return domain.AdminTrend{}, false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	return s.repo.Get(ctx, slug)
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
