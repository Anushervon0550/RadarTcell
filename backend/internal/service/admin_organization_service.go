package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminOrganizationService struct {
	repo  ports.AdminOrganizationRepository
	cache ports.Cache
}

func NewAdminOrganizationService(repo ports.AdminOrganizationRepository, cache ports.Cache) *AdminOrganizationService {
	return &AdminOrganizationService{repo: repo, cache: cache}
}

func (s *AdminOrganizationService) Create(ctx context.Context, cmd domain.OrganizationUpsert) (string, error) {
	if err := validateOrg(cmd, true); err != nil {
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

func (s *AdminOrganizationService) Update(ctx context.Context, slug string, cmd domain.OrganizationUpsert) (string, bool, error) {
	if err := validateSlugValue(slug); err != nil {
		return "", false, err
	}
	if err := validateOrg(cmd, false); err != nil {
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

func (s *AdminOrganizationService) Delete(ctx context.Context, slug string) (bool, error) {
	if err := validateSlugValue(slug); err != nil {
		return false, err
	}
	ok, err := s.repo.Delete(ctx, slug)
	if err != nil || !ok {
		return ok, err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return ok, nil
}

func (s *AdminOrganizationService) List(ctx context.Context) ([]domain.Organization, error) {
	return s.repo.List(ctx)
}

func (s *AdminOrganizationService) Get(ctx context.Context, slug string) (domain.Organization, bool, error) {
	if err := validateSlugValue(slug); err != nil {
		return domain.Organization{}, false, err
	}
	return s.repo.Get(ctx, slug)
}

func validateOrg(cmd domain.OrganizationUpsert, isCreate bool) error {
	if isCreate {
		if err := validateSlugValue(cmd.Slug); err != nil {
			return err
		}
	}
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	return nil
}
