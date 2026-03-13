package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminTagService struct {
	repo  ports.AdminTagRepository
	cache ports.Cache
}

func NewAdminTagService(repo ports.AdminTagRepository, cache ports.Cache) *AdminTagService {
	return &AdminTagService{repo: repo, cache: cache}
}

func (s *AdminTagService) Create(ctx context.Context, cmd domain.TagUpsert) (string, error) {
	if err := validateTag(cmd, true); err != nil {
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

func (s *AdminTagService) Update(ctx context.Context, slug string, cmd domain.TagUpsert) (string, bool, error) {
	if err := validateSlugValue(slug); err != nil {
		return "", false, err
	}
	if err := validateTag(cmd, false); err != nil {
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

func (s *AdminTagService) Delete(ctx context.Context, slug string) (bool, error) {
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

func (s *AdminTagService) List(ctx context.Context) ([]domain.Tag, error) {
	return s.repo.List(ctx)
}

func (s *AdminTagService) Get(ctx context.Context, slug string) (domain.Tag, bool, error) {
	if err := validateSlugValue(slug); err != nil {
		return domain.Tag{}, false, err
	}
	return s.repo.Get(ctx, slug)
}

func validateTag(cmd domain.TagUpsert, isCreate bool) error {
	if isCreate {
		if err := validateSlugValue(cmd.Slug); err != nil {
			return err
		}
	}
	if strings.TrimSpace(cmd.Title) == "" {
		return fmt.Errorf("%w: title is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Category) == "" {
		return fmt.Errorf("%w: category is required", domain.ErrInvalid)
	}
	return nil
}
