package warehouses_repository_postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
)

type WarehousesRepository struct {
	tx *core_postgres_tx.Tx
}

func NewWarehousesRepository(tx *core_postgres_tx.Tx) *WarehousesRepository {
	return &WarehousesRepository{
		tx: tx,
	}
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
