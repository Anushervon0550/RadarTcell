package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminTagService struct {
	repo ports.AdminTagRepository
}

func NewAdminTagService(repo ports.AdminTagRepository) *AdminTagService {
	return &AdminTagService{repo: repo}
}

func (s *AdminTagService) Create(ctx context.Context, cmd domain.TagUpsert) (string, error) {
	if err := validateTag(cmd, true); err != nil {
		return "", err
	}
	return s.repo.Create(ctx, cmd)
}

func (s *AdminTagService) Update(ctx context.Context, slug string, cmd domain.TagUpsert) (string, bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if err := validateTag(cmd, false); err != nil {
		return "", false, err
	}
	return s.repo.Update(ctx, slug, cmd)
}

func (s *AdminTagService) Delete(ctx context.Context, slug string) (bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	return s.repo.Delete(ctx, slug)
}

func validateTag(cmd domain.TagUpsert, isCreate bool) error {
	if isCreate && strings.TrimSpace(cmd.Slug) == "" {
		return fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Title) == "" {
		return fmt.Errorf("%w: title is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Category) == "" {
		return fmt.Errorf("%w: category is required", domain.ErrInvalid)
	}
	return nil
}
