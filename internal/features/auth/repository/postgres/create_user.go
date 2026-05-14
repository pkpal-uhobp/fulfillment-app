package auth_repository_postgres

import (
	"context"
	"database/sql"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *AuthRepository) CreateUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		INSERT INTO users (
			email,
			password_hash,
			full_name,
			phone,
			role
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, password_hash, full_name, phone, role, is_active, is_blocked;
	`

	var (
		result    core_domain.User
		phone     sql.NullString
		roleValue string
	)

	err := q.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Phone,
		user.Role.String(),
	).Scan(
		&result.ID,
		&result.Email,
		&result.PasswordHash,
		&result.FullName,
		&phone,
		&roleValue,
		&result.IsActive,
		&result.IsBlocked,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return core_domain.User{}, fmt.Errorf(
				"%w: user already exists",
				core_errors.ErrConflict,
			)
		}

		return core_domain.User{}, fmt.Errorf("create user: %w", err)
	}

	if phone.Valid {
		result.Phone = &phone.String
	}
	result.Role = core_domain.Role(roleValue)

	return result, nil
}
