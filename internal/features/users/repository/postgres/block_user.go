package users_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *UsersRepository) BlockUser(
	ctx context.Context,
	userID int64,
	reason *string,
) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)
	const query = `
		UPDATE users
		SET is_blocked = TRUE,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1;
	`

	result, err := q.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("block user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: user not found", core_errors.ErrNotFound)
	}
	return nil
}

func (r *UsersRepository) DeactivateUser(
	ctx context.Context,
	userID int64,
) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)
	const query = `
		UPDATE users
		SET is_active = FALSE,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1;
	`

	result, err := q.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("deactivate user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: user not found", core_errors.ErrNotFound)
	}
	return nil
}
