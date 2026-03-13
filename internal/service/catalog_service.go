package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/google/uuid"
)

type CatalogService struct {
	repo  ports.CatalogRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewCatalogService(repo ports.CatalogRepository, cache ports.Cache, ttl time.Duration) *CatalogService {
	return &CatalogService{repo: repo, cache: cache, ttl: ttl}
}

func (s *CatalogService) ListTrends(ctx context.Context, locale string) ([]domain.Trend, error) {
	var cached []domain.Trend
	suffix := s.cacheSuffix("catalog:trends", locale)
	if s.getCachedList(ctx, suffix, &cached) {
		return cached, nil
	}
	items, err := s.repo.ListTrends(ctx, locale)
	if err != nil {
		return nil, err
	}
	s.setCachedList(ctx, suffix, items)
	return items, nil
}

func (s *CatalogService) ListSDGs(ctx context.Context, locale string) ([]domain.SDG, error) {
	var cached []domain.SDG
	suffix := s.cacheSuffix("catalog:sdgs", locale)
	if s.getCachedList(ctx, suffix, &cached) {
		return cached, nil
	}
	items, err := s.repo.ListSDGs(ctx, locale)
	if err != nil {
		return nil, err
	}
	s.setCachedList(ctx, suffix, items)
	return items, nil
}

func (s *CatalogService) ListTags(ctx context.Context, locale string) ([]domain.Tag, error) {
	var cached []domain.Tag
	suffix := s.cacheSuffix("catalog:tags", locale)
	if s.getCachedList(ctx, suffix, &cached) {
		return cached, nil
	}
	items, err := s.repo.ListTags(ctx, locale)
	if err != nil {
		return nil, err
	}
	s.setCachedList(ctx, suffix, items)
	return items, nil
}

func (s *CatalogService) ListOrganizations(ctx context.Context, locale string) ([]domain.Organization, error) {
	var cached []domain.Organization
	suffix := s.cacheSuffix("catalog:organizations", locale)
	if s.getCachedList(ctx, suffix, &cached) {
		return cached, nil
	}
	items, err := s.repo.ListOrganizations(ctx, locale)
	if err != nil {
		return nil, err
	}
	s.setCachedList(ctx, suffix, items)
	return items, nil
}

func (s *CatalogService) ListMetrics(ctx context.Context, locale string) ([]domain.MetricDefinition, error) {
	var cached []domain.MetricDefinition
	suffix := s.cacheSuffix("catalog:metrics", locale)
	if s.getCachedList(ctx, suffix, &cached) {
		return cached, nil
	}
	items, err := s.repo.ListMetrics(ctx, locale)
	if err != nil {
		return nil, err
	}
	s.setCachedList(ctx, suffix, items)
	return items, nil
}

func (s *CatalogService) GetOrganizationBySlug(ctx context.Context, slug string) (domain.Organization, bool, error) {
	return s.repo.GetOrganizationBySlug(ctx, slug)
}
func (s *CatalogService) GetMetricValue(ctx context.Context, metricID, technologyID string) (map[string]any, bool, error) {
	metricID = strings.TrimSpace(metricID)
	technologyID = strings.TrimSpace(technologyID)

	if metricID == "" {
		return nil, false, fmt.Errorf("%w: metric id is required", domain.ErrInvalid)
	}
	if technologyID == "" {
		return nil, false, fmt.Errorf("%w: technology_id is required", domain.ErrInvalid)
	}
	if _, err := uuid.Parse(metricID); err != nil {
		return nil, false, fmt.Errorf("%w: metric id must be uuid", domain.ErrInvalid)
	}
	if _, err := uuid.Parse(technologyID); err != nil {
		return nil, false, fmt.Errorf("%w: technology_id must be uuid", domain.ErrInvalid)
	}

	return s.repo.GetMetricValue(ctx, metricID, technologyID)
}

func (s *CatalogService) getCachedList(ctx context.Context, suffix string, out any) bool {
	if s.cache == nil || s.ttl <= 0 {
		return false
	}
	key := s.cacheKey(ctx, suffix)
	b, ok, err := s.cache.Get(ctx, key)
	if err != nil || !ok {
		return false
	}
	if err := json.Unmarshal(b, out); err != nil {
		return false
	}
	return true
}

func (s *CatalogService) setCachedList(ctx context.Context, suffix string, v any) {
	if s.cache == nil || s.ttl <= 0 {
		return
	}
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	_ = s.cache.Set(ctx, s.cacheKey(ctx, suffix), b, s.ttl)
}

func (s *CatalogService) cacheKey(ctx context.Context, suffix string) string {
	v := cacheVersion(ctx, s.cache, cacheVersionCatalog)
	return fmt.Sprintf("%s:%s:%s", cacheVersionCatalog, v, suffix)
}

func (s *CatalogService) cacheSuffix(base, locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if locale == "" {
		return base
	}
	return base + ":" + locale
}
