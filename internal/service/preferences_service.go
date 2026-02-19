package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type PreferencesService struct {
	repo ports.PreferencesRepository
}

func NewPreferencesService(repo ports.PreferencesRepository) *PreferencesService {
	return &PreferencesService{repo: repo}
}

func (s *PreferencesService) Save(ctx context.Context, p domain.Preferences) error {
	p.UserID = strings.TrimSpace(p.UserID)
	if len(p.Settings) == 0 {
		p.Settings = json.RawMessage(`{}`)
	}
	return s.repo.UpsertPreferences(ctx, p.UserID, p.Settings)
}

func (s *PreferencesService) Get(ctx context.Context, userID string) (domain.Preferences, bool, error) {
	userID = strings.TrimSpace(userID)
	settings, ok, err := s.repo.GetPreferences(ctx, userID)
	if err != nil || !ok {
		return domain.Preferences{}, ok, err
	}
	return domain.Preferences{UserID: userID, Settings: settings}, true, nil
}
