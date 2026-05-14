package auth_repository_postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *AuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (core_domain.User, error) {
	q := r.tx.Querier(ctx)

	const query = `
		SELECT
			id,
			email,
			password_hash,
			full_name,
			phone,
			role,
			is_active,
			is_blocked
		FROM users
		WHERE lower(email) = lower($1);
	`

	return scanUser(q.QueryRow(ctx, query, email))
}

func (r *AuthRepository) GetUserByID(
	ctx context.Context,
	userID int64,
) (core_domain.User, error) {
	q := r.tx.Querier(ctx)

	const query = `
		SELECT
			id,
			email,
			password_hash,
			full_name,
			phone,
			role,
			is_active,
			is_blocked
		FROM users
		WHERE id = $1;
	`

	return scanUser(q.QueryRow(ctx, query, userID))
}

func scanUser(row pgx.Row) (core_domain.User, error) {
	var (
		user      core_domain.User
		phone     sql.NullString
		roleValue string
	)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&phone,
		&roleValue,
		&user.IsActive,
		&user.IsBlocked,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf(
				"%w: user not found",
				core_errors.ErrNotFound,
			)
		}

		return core_domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	if phone.Valid {
		user.Phone = &phone.String
	}

	user.Role = core_domain.Role(roleValue)

	return user, nil
}
