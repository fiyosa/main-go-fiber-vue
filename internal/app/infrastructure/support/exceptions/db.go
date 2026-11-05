package exceptions

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// PostgreSQL SQLSTATE codes handled explicitly.
const (
	sqlStateUniqueViolation     = "23505"
	sqlStateForeignKeyViolation = "23503"
	sqlStateNotNullViolation    = "23502"
	sqlStateSerialization       = "40001"
	sqlStateDeadlock            = "40P01"
)

// FromDB maps GORM and PostgreSQL errors into a typed exception. It returns
// nil when err is nil so it can wrap repository return values directly.
func FromDB(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case sqlStateUniqueViolation:
			return New("db.duplicate", http.StatusConflict, "db.duplicate", err)
		case sqlStateForeignKeyViolation, sqlStateNotNullViolation:
			return New("db.constraint", http.StatusUnprocessableEntity, "db.constraint_failed", err)
		case sqlStateSerialization, sqlStateDeadlock:
			return ServiceUnavailable(err)
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return New("db.timeout", http.StatusGatewayTimeout, "common.try_again_later", err)
	}
	if errors.Is(err, context.Canceled) {
		return ServiceUnavailable(err)
	}

	return New("db.internal", http.StatusInternalServerError, "common.try_again_later", err)
}
