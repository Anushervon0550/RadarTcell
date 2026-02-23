package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminTechnologyService struct {
	repo ports.AdminTechnologyRepository
}

func NewAdminTechnologyService(repo ports.AdminTechnologyRepository) *AdminTechnologyService {
	return &AdminTechnologyService{repo: repo}
}

func (s *AdminTechnologyService) Create(ctx context.Context, cmd domain.TechnologyUpsert) (string, error) {
	if err := validateTechUpsert(cmd, true); err != nil {
		return "", err
	}
	return s.repo.Create(ctx, cmd)
}

func (s *AdminTechnologyService) Update(ctx context.Context, slug string, cmd domain.TechnologyUpsert) (string, bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if err := validateTechUpsert(cmd, false); err != nil {
		return "", false, err
	}
	return s.repo.Update(ctx, slug, cmd)
}

func (s *AdminTechnologyService) Delete(ctx context.Context, slug string) (bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	return s.repo.Delete(ctx, slug)
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
