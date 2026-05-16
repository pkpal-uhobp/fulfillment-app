package users_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)
	query := fmt.Sprintf(`
		INSERT INTO users (
			email,
			password_hash,
			full_name,
			phone,
			role
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING %s;
	`, userColumns)

	created, err := scanUserRow(q.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Phone,
		user.Role.String(),
	))
	if err != nil {
		if isUniqueViolation(err) {
			return core_domain.User{}, fmt.Errorf("%w: user already exists", core_errors.ErrConflict)
		}
		return core_domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return created, nil
}
