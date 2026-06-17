package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PreferencesRepo struct {
	db *pgxpool.Pool
}

func NewPreferencesRepo(db *pgxpool.Pool) *PreferencesRepo {
	return &PreferencesRepo{db: db}
}

func (r *PreferencesRepo) UpsertPreferences(ctx context.Context, userID string, settings json.RawMessage) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_preferences (user_id, settings)
		VALUES ($1, $2::jsonb)
		ON CONFLICT (user_id)
		DO UPDATE SET settings = EXCLUDED.settings, updated_at = now()
	`, userID, settings)
	if err != nil {
		return fmt.Errorf("upsert preferences: %w", err)
	}
	return nil
}

func (r *PreferencesRepo) GetPreferences(ctx context.Context, userID string) (json.RawMessage, bool, error) {
	var settings json.RawMessage
	err := r.db.QueryRow(ctx, `
		SELECT settings
		FROM user_preferences
		WHERE user_id = $1
	`, userID).Scan(&settings)

	if err == nil {
		return settings, true, nil
	}
	if err == pgx.ErrNoRows {
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("get preferences: %w", err)
}
