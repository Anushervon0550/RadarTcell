package postgres

import (
	"errors"
	"fmt"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

// mapPGErr translates Postgres errors into domain errors.
// 23505 = unique_violation, 23503 = foreign_key_violation -> ErrConflict.
// All other DB errors are returned as a generic wrapped error.
func mapPGErr(err error, conflictMsg string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505", "23503":
			return fmt.Errorf("%w: %s", domain.ErrConflict, conflictMsg)
		}
	}
	return fmt.Errorf("db error: %w", err)
}
