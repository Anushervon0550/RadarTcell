package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type AdminSDGService struct {
	repo ports.AdminSDGRepository
}

func NewAdminSDGService(repo ports.AdminSDGRepository) *AdminSDGService {
	return &AdminSDGService{repo: repo}
}

func (s *AdminSDGService) Create(ctx context.Context, cmd domain.SDGUpsert) (string, error) {
	if err := validateSDG(&cmd, true); err != nil {
		return "", err
	}
	return s.repo.Create(ctx, cmd)
}

func (s *AdminSDGService) Update(ctx context.Context, code string, cmd domain.SDGUpsert) (bool, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return false, fmt.Errorf("%w: code is required", domain.ErrInvalid)
	}
	// code меняем через path, в body он не обязателен
	if err := validateSDG(&cmd, false); err != nil {
		return false, err
	}
	return s.repo.Update(ctx, code, cmd)
}

func (s *AdminSDGService) Delete(ctx context.Context, code string) (bool, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return false, fmt.Errorf("%w: code is required", domain.ErrInvalid)
	}
	return s.repo.Delete(ctx, code)
}

func validateSDG(cmd *domain.SDGUpsert, requireCode bool) error {
	if requireCode && strings.TrimSpace(cmd.Code) == "" {
		return fmt.Errorf("%w: code is required", domain.ErrInvalid)
	}
	if strings.TrimSpace(cmd.Title) == "" {
		return fmt.Errorf("%w: title is required", domain.ErrInvalid)
	}
	cmd.Code = strings.TrimSpace(cmd.Code)
	cmd.Title = strings.TrimSpace(cmd.Title)
	return nil
}
