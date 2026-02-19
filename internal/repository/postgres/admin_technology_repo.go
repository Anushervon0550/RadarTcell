package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminTechnologyRepo struct {
	db *pgxpool.Pool
}

func NewAdminTechnologyRepo(db *pgxpool.Pool) *AdminTechnologyRepo {
	return &AdminTechnologyRepo{db: db}
}

var _ ports.AdminTechnologyRepository = (*AdminTechnologyRepo)(nil)

func (r *AdminTechnologyRepo) Create(ctx context.Context, cmd domain.TechnologyUpsert) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	trendID, err := resolveTrendID(ctx, tx, cmd.TrendSlug)
	if err != nil {
		return "", err
	}

	var techID string
	err = tx.QueryRow(ctx, `
		INSERT INTO technologies (
			slug, list_index, name,
			description_short, description_full,
			readiness_level,
			custom_metric_1, custom_metric_2, custom_metric_3, custom_metric_4,
			image_url, source_link,
			trend_id
		)
		VALUES (
			$1, $2, $3,
			$4, $5,
			$6,
			$7, $8, $9, $10,
			$11, $12,
			$13::uuid
		)
		RETURNING id::text
	`,
		strings.TrimSpace(cmd.Slug),
		cmd.Index,
		strings.TrimSpace(cmd.Name),
		cmd.DescriptionShort,
		cmd.DescriptionFull,
		cmd.TRL,
		cmd.CustomMetric1, cmd.CustomMetric2, cmd.CustomMetric3, cmd.CustomMetric4,
		cmd.ImageURL,
		cmd.SourceLink,
		trendID,
	).Scan(&techID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("%w: technology slug already exists", domain.ErrConflict)
		}
		return "", fmt.Errorf("insert technology: %w", err)
	}

	if err := replaceLinks(ctx, tx, techID, cmd); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}
	return techID, nil
}

func (r *AdminTechnologyRepo) Update(ctx context.Context, slug string, cmd domain.TechnologyUpsert) (string, bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var techID string
	err = tx.QueryRow(ctx, `SELECT id::text FROM technologies WHERE slug=$1`, slug).Scan(&techID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("find technology: %w", err)
	}

	trendID, err := resolveTrendID(ctx, tx, cmd.TrendSlug)
	if err != nil {
		return "", false, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE technologies
		SET
			list_index=$2,
			name=$3,
			description_short=$4,
			description_full=$5,
			readiness_level=$6,
			custom_metric_1=$7,
			custom_metric_2=$8,
			custom_metric_3=$9,
			custom_metric_4=$10,
			image_url=$11,
			source_link=$12,
			trend_id=$13::uuid,
			updated_at=now()
		WHERE id=$1::uuid
	`,
		techID,
		cmd.Index,
		strings.TrimSpace(cmd.Name),
		cmd.DescriptionShort,
		cmd.DescriptionFull,
		cmd.TRL,
		cmd.CustomMetric1, cmd.CustomMetric2, cmd.CustomMetric3, cmd.CustomMetric4,
		cmd.ImageURL,
		cmd.SourceLink,
		trendID,
	)
	if err != nil {
		return "", false, fmt.Errorf("update technology: %w", err)
	}

	if err := replaceLinks(ctx, tx, techID, cmd); err != nil {
		return "", false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", false, fmt.Errorf("commit: %w", err)
	}
	return techID, true, nil
}

func (r *AdminTechnologyRepo) Delete(ctx context.Context, slug string) (bool, error) {
	ct, err := r.db.Exec(ctx, `DELETE FROM technologies WHERE slug=$1`, slug)
	if err != nil {
		return false, fmt.Errorf("delete technology: %w", err)
	}
	return ct.RowsAffected() > 0, nil
}

// ---- helpers ----

func resolveTrendID(ctx context.Context, q pgx.Tx, trendSlug string) (string, error) {
	trendSlug = strings.TrimSpace(trendSlug)
	var id string
	err := q.QueryRow(ctx, `SELECT id::text FROM trends WHERE slug=$1`, trendSlug).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("%w: trend not found", domain.ErrInvalid)
	}
	return "", fmt.Errorf("resolve trend: %w", err)
}

func replaceLinks(ctx context.Context, tx pgx.Tx, techID string, cmd domain.TechnologyUpsert) error {
	// очистить старые связи
	if _, err := tx.Exec(ctx, `DELETE FROM technology_tags WHERE technology_id=$1::uuid`, techID); err != nil {
		return fmt.Errorf("clear tech tags: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM technology_sdgs WHERE technology_id=$1::uuid`, techID); err != nil {
		return fmt.Errorf("clear tech sdgs: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM technology_organizations WHERE technology_id=$1::uuid`, techID); err != nil {
		return fmt.Errorf("clear tech orgs: %w", err)
	}

	// tags by slug
	if len(cmd.TagSlugs) > 0 {
		tagIDs, err := resolveIDsByText(ctx, tx,
			`SELECT id::text, slug FROM tags WHERE slug = ANY($1::text[])`,
			cmd.TagSlugs,
		)
		if err != nil {
			return err
		}
		for _, id := range tagIDs {
			if _, err := tx.Exec(ctx, `
				INSERT INTO technology_tags (technology_id, tag_id)
				VALUES ($1::uuid, $2::uuid)
				ON CONFLICT DO NOTHING
			`, techID, id); err != nil {
				return fmt.Errorf("insert tech tag: %w", err)
			}
		}
	}

	// sdgs by code
	if len(cmd.SDGCodes) > 0 {
		sdgIDs, err := resolveIDsByText(ctx, tx,
			`SELECT id::text, code FROM sdgs WHERE code = ANY($1::text[])`,
			cmd.SDGCodes,
		)
		if err != nil {
			return err
		}
		for _, id := range sdgIDs {
			if _, err := tx.Exec(ctx, `
				INSERT INTO technology_sdgs (technology_id, sdg_id)
				VALUES ($1::uuid, $2::uuid)
				ON CONFLICT DO NOTHING
			`, techID, id); err != nil {
				return fmt.Errorf("insert tech sdg: %w", err)
			}
		}
	}

	// orgs by slug
	if len(cmd.OrganizationSlugs) > 0 {
		orgIDs, err := resolveIDsByText(ctx, tx,
			`SELECT id::text, slug FROM organizations WHERE slug = ANY($1::text[])`,
			cmd.OrganizationSlugs,
		)
		if err != nil {
			return err
		}
		for _, id := range orgIDs {
			if _, err := tx.Exec(ctx, `
				INSERT INTO technology_organizations (technology_id, organization_id)
				VALUES ($1::uuid, $2::uuid)
				ON CONFLICT DO NOTHING
			`, techID, id); err != nil {
				return fmt.Errorf("insert tech org: %w", err)
			}
		}
	}

	return nil
}

// Возвращает ID-шники и валидирует, что все requested значения существуют.
func resolveIDsByText(ctx context.Context, tx pgx.Tx, sql string, wanted []string) ([]string, error) {
	wantedSet := map[string]struct{}{}
	for _, v := range wanted {
		v = strings.TrimSpace(v)
		if v != "" {
			wantedSet[v] = struct{}{}
		}
	}
	if len(wantedSet) == 0 {
		return nil, nil
	}

	uniq := make([]string, 0, len(wantedSet))
	for k := range wantedSet {
		uniq = append(uniq, k)
	}

	rows, err := tx.Query(ctx, sql, uniq)
	if err != nil {
		return nil, fmt.Errorf("resolve ids query: %w", err)
	}
	defer rows.Close()

	found := map[string]string{} // key -> id
	for rows.Next() {
		var id, key string
		if err := rows.Scan(&id, &key); err != nil {
			return nil, fmt.Errorf("resolve ids scan: %w", err)
		}
		found[key] = id
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("resolve ids rows: %w", err)
	}

	// validate missing
	for k := range wantedSet {
		if _, ok := found[k]; !ok {
			return nil, fmt.Errorf("%w: reference not found: %s", domain.ErrInvalid, k)
		}
	}

	out := make([]string, 0, len(found))
	for _, id := range found {
		out = append(out, id)
	}
	return out, nil
}
