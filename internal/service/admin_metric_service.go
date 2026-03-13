package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

var metricFieldKeyRe = regexp.MustCompile(`^[a-z][a-z0-9_]{1,62}$`)

type AdminMetricService struct {
	repo  ports.AdminMetricRepository
	cache ports.Cache
}

func NewAdminMetricService(repo ports.AdminMetricRepository, cache ports.Cache) *AdminMetricService {
	return &AdminMetricService{repo: repo, cache: cache}
}

func (s *AdminMetricService) Create(ctx context.Context, cmd domain.MetricDefinitionUpsert) (string, error) {
	if err := validateMetric(&cmd); err != nil {
		return "", err
	}
	fk, err := normalizeAndValidateMetricFieldKey(cmd.FieldKey)
	if err != nil {
		return "", err
	}
	cmd.FieldKey = fk
	id, err := s.repo.Create(ctx, cmd)
	if err != nil {
		return "", err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return id, nil
}

func (s *AdminMetricService) Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return false, fmt.Errorf("%w: id is required", domain.ErrInvalid)
	}
	if err := validateMetric(&cmd); err != nil {
		return false, err
	}
	fk, err := normalizeAndValidateMetricFieldKey(cmd.FieldKey)
	if err != nil {
		return false, err
	}
	cmd.FieldKey = fk
	ok, err := s.repo.Update(ctx, id, cmd)
	if err != nil || !ok {
		return ok, err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return ok, nil
}

func (s *AdminMetricService) Delete(ctx context.Context, id string) (bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return false, fmt.Errorf("%w: id is required", domain.ErrInvalid)
	}
	ok, err := s.repo.Delete(ctx, id)
	if err != nil || !ok {
		return ok, err
	}
	bumpCacheVersion(ctx, s.cache, cacheVersionCatalog)
	bumpCacheVersion(ctx, s.cache, cacheVersionTechnologies)
	return ok, nil
}

func (s *AdminMetricService) List(ctx context.Context) ([]domain.MetricDefinition, error) {
	return s.repo.List(ctx)
}

func (s *AdminMetricService) Get(ctx context.Context, id string) (domain.MetricDefinition, bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.MetricDefinition{}, false, fmt.Errorf("%w: id is required", domain.ErrInvalid)
	}
	return s.repo.Get(ctx, id)
}

func validateMetric(cmd *domain.MetricDefinitionUpsert) error {
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}

	t := strings.TrimSpace(cmd.Type)
	switch t {
	case "bubble", "bar", "distance":
		// ok
	default:
		return fmt.Errorf("%w: type must be bubble|bar|distance", domain.ErrInvalid)
	}

	cmd.Type = t
	cmd.Name = strings.TrimSpace(cmd.Name)
	return nil
}
func normalizeAndValidateMetricFieldKey(v *string) (*string, error) {
	if v == nil {
		return nil, nil
	}

	s := strings.TrimSpace(*v)
	if s == "" {
		return nil, nil
	}

	if !metricFieldKeyRe.MatchString(s) {
		return nil, fmt.Errorf("%w: field_key must match ^[a-z][a-z0-9_]{1,62}$", domain.ErrInvalid)
	}

	return &s, nil
}
