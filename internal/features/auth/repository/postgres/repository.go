package auth_repository_postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
)

type AuthRepository struct {
	tx *core_postgres_tx.Tx
}

func NewAuthRepository(tx *core_postgres_tx.Tx) *AuthRepository {
	return &AuthRepository{
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
