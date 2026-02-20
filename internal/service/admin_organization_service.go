package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminOrganizationService struct {
	repo ports.AdminOrganizationRepository
}

func NewAdminOrganizationService(repo ports.AdminOrganizationRepository) *AdminOrganizationService {
	return &AdminOrganizationService{repo: repo}
}

func (s *AdminOrganizationService) Create(ctx context.Context, cmd domain.OrganizationUpsert) (string, error) {
	if err := validateOrg(cmd, true); err != nil {
		return "", err
	}
	return s.repo.Create(ctx, cmd)
}

func (s *AdminOrganizationService) Update(ctx context.Context, slug string, cmd domain.OrganizationUpsert) (string, bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if err := validateOrg(cmd, false); err != nil {
		return "", false, err
	}
	return s.repo.Update(ctx, slug, cmd)
}

func (s *AdminOrganizationService) Delete(ctx context.Context, slug string) (bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return false, fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	return s.repo.Delete(ctx, slug)
}

func validateOrg(cmd domain.OrganizationUpsert, isCreate bool) error {
	if isCreate && strings.TrimSpace(cmd.Slug) == "" {
		return fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	return nil
}
