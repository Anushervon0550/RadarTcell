package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminTechnologyService struct {
	repo  ports.AdminTechnologyRepository
	cache ports.Cache
}

func NewAdminTechnologyService(repo ports.AdminTechnologyRepository, cache ports.Cache) *AdminTechnologyService {
	return &AdminTechnologyService{repo: repo, cache: cache}
}

func (s *AdminTechnologyService) Create(ctx context.Context, cmd domain.TechnologyUpsert) (string, error) {
	if err := validateTechUpsert(cmd, true); err != nil {
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

func (s *AdminTechnologyService) Update(ctx context.Context, slug string, cmd domain.TechnologyUpsert) (string, bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if err := validateTechUpsert(cmd, false); err != nil {
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

func (s *AdminTechnologyService) Delete(ctx context.Context, slug string) (bool, error) {
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

func (s *AdminTechnologyService) List(ctx context.Context, p domain.AdminTechnologyListParams) (domain.AdminTechnologyListResult, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 50
	}
	if p.Limit > 200 {
		return domain.AdminTechnologyListResult{}, fmt.Errorf("%w: limit must be between 1 and 200", domain.ErrInvalid)
	}

	items, total, err := s.repo.List(ctx, p)
	if err != nil {
		return domain.AdminTechnologyListResult{}, err
	}

	return domain.AdminTechnologyListResult{
		Page:  p.Page,
		Limit: p.Limit,
		Total: total,
		Items: items,
	}, nil
}

func (s *AdminTechnologyService) Get(ctx context.Context, slug string) (domain.TechnologyAdmin, bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return domain.TechnologyAdmin{}, false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	return s.repo.Get(ctx, slug)
}

func validateTechUpsert(cmd domain.TechnologyUpsert, isCreate bool) error {
	if isCreate {
		if strings.TrimSpace(cmd.Slug) == "" {
			return fmt.Errorf("%w: slug is required", domain.ErrInvalid)
		}
	}
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.TrendSlug) == "" {
		return fmt.Errorf("%w: trend_slug is required", domain.ErrInvalid)
	}
	if cmd.TRL < 1 || cmd.TRL > 9 {
		return fmt.Errorf("%w: trl must be 1..9", domain.ErrInvalid)
	}

	if cmd.Index < 1 || cmd.Index > 99 {
		return fmt.Errorf("%w: index must be between 1 and 99", domain.ErrInvalid)
	}
	if cmd.TRL < 1 || cmd.TRL > 9 {
		return fmt.Errorf("%w: trl must be between 1 and 9", domain.ErrInvalid)
	}

	return nil
}
