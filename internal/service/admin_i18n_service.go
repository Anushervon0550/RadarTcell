package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminI18nService struct {
	repo  ports.AdminI18nRepository
	cache ports.Cache
}

func NewAdminI18nService(repo ports.AdminI18nRepository, cache ports.Cache) *AdminI18nService {
	return &AdminI18nService{repo: repo, cache: cache}
}

func (s *AdminI18nService) UpsertTrend(ctx context.Context, trendSlug string, cmd domain.TrendI18nUpsert) error {
	if err := validateLocale(cmd.Locale); err != nil {
		return err
	}
	if strings.TrimSpace(trendSlug) == "" {
		return fmt.Errorf("%w: trend_slug is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	if err := s.repo.UpsertTrend(ctx, trendSlug, cmd); err != nil {
		return err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return nil
}

func (s *AdminI18nService) GetTrend(ctx context.Context, trendSlug, locale string) (domain.TrendI18n, bool, error) {
	if err := validateLocale(locale); err != nil {
		return domain.TrendI18n{}, false, err
	}
	if strings.TrimSpace(trendSlug) == "" {
		return domain.TrendI18n{}, false, fmt.Errorf("%w: trend_slug is required", domain.ErrInvalid)
	}
	return s.repo.GetTrend(ctx, trendSlug, locale)
}

func (s *AdminI18nService) DeleteTrend(ctx context.Context, trendSlug, locale string) (bool, error) {
	if err := validateLocale(locale); err != nil {
		return false, err
	}
	if strings.TrimSpace(trendSlug) == "" {
		return false, fmt.Errorf("%w: trend_slug is required", domain.ErrInvalid)
	}
	ok, err := s.repo.DeleteTrend(ctx, trendSlug, locale)
	if err != nil || !ok {
		return ok, err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return ok, nil
}

func (s *AdminI18nService) UpsertTechnology(ctx context.Context, techSlug string, cmd domain.TechnologyI18nUpsert) error {
	if err := validateLocale(cmd.Locale); err != nil {
		return err
	}
	if strings.TrimSpace(techSlug) == "" {
		return fmt.Errorf("%w: tech_slug is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	if err := s.repo.UpsertTechnology(ctx, techSlug, cmd); err != nil {
		return err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return nil
}

func (s *AdminI18nService) GetTechnology(ctx context.Context, techSlug, locale string) (domain.TechnologyI18n, bool, error) {
	if err := validateLocale(locale); err != nil {
		return domain.TechnologyI18n{}, false, err
	}
	if strings.TrimSpace(techSlug) == "" {
		return domain.TechnologyI18n{}, false, fmt.Errorf("%w: tech_slug is required", domain.ErrInvalid)
	}
	return s.repo.GetTechnology(ctx, techSlug, locale)
}

func (s *AdminI18nService) DeleteTechnology(ctx context.Context, techSlug, locale string) (bool, error) {
	if err := validateLocale(locale); err != nil {
		return false, err
	}
	if strings.TrimSpace(techSlug) == "" {
		return false, fmt.Errorf("%w: tech_slug is required", domain.ErrInvalid)
	}
	ok, err := s.repo.DeleteTechnology(ctx, techSlug, locale)
	if err != nil || !ok {
		return ok, err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return ok, nil
}

func (s *AdminI18nService) UpsertMetric(ctx context.Context, metricID string, cmd domain.MetricI18nUpsert) error {
	if err := validateLocale(cmd.Locale); err != nil {
		return err
	}
	metricID = strings.TrimSpace(metricID)
	if metricID == "" {
		return fmt.Errorf("%w: metric_id is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	if err := s.repo.UpsertMetric(ctx, metricID, cmd); err != nil {
		return err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return nil
}

func (s *AdminI18nService) GetMetric(ctx context.Context, metricID, locale string) (domain.MetricI18n, bool, error) {
	if err := validateLocale(locale); err != nil {
		return domain.MetricI18n{}, false, err
	}
	metricID = strings.TrimSpace(metricID)
	if metricID == "" {
		return domain.MetricI18n{}, false, fmt.Errorf("%w: metric_id is required", domain.ErrInvalid)
	}
	return s.repo.GetMetric(ctx, metricID, locale)
}

func (s *AdminI18nService) DeleteMetric(ctx context.Context, metricID, locale string) (bool, error) {
	if err := validateLocale(locale); err != nil {
		return false, err
	}
	metricID = strings.TrimSpace(metricID)
	if metricID == "" {
		return false, fmt.Errorf("%w: metric_id is required", domain.ErrInvalid)
	}
	ok, err := s.repo.DeleteMetric(ctx, metricID, locale)
	if err != nil || !ok {
		return ok, err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return ok, nil
}

func validateLocale(locale string) error {
	locale = strings.TrimSpace(locale)
	if locale == "" {
		return fmt.Errorf("%w: locale is required", domain.ErrInvalid)
	}
	return nil
}
