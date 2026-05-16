package pickupcalendar_repository_postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
)

type PickupCalendarRepository struct {
	tx *core_postgres_tx.Tx
}

func NewPickupCalendarRepository(tx *core_postgres_tx.Tx) *PickupCalendarRepository {
	return &PickupCalendarRepository{tx: tx}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}

func isCheckViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23514"
	}
	return false
}
