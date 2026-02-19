package service

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type CatalogService struct {
	repo ports.CatalogRepository
}

func NewCatalogService(repo ports.CatalogRepository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) ListTrends(ctx context.Context) ([]domain.Trend, error) {
	return s.repo.ListTrends(ctx)
}
func (s *CatalogService) ListSDGs(ctx context.Context) ([]domain.SDG, error) {
	return s.repo.ListSDGs(ctx)
}
func (s *CatalogService) ListTags(ctx context.Context) ([]domain.Tag, error) {
	return s.repo.ListTags(ctx)
}
func (s *CatalogService) ListOrganizations(ctx context.Context) ([]domain.Organization, error) {
	return s.repo.ListOrganizations(ctx)
}
func (s *CatalogService) ListMetrics(ctx context.Context) ([]domain.MetricDefinition, error) {
	return s.repo.ListMetrics(ctx)
}
func (s *CatalogService) GetOrganizationBySlug(ctx context.Context, slug string) (domain.Organization, bool, error) {
	return s.repo.GetOrganizationBySlug(ctx, slug)
}
