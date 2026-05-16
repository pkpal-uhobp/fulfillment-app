package users_repository_postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

const userColumns = `
	id,
	email,
	password_hash,
	full_name,
	phone,
	role,
	is_active,
	is_blocked
`

func scanUserRow(row pgx.Row) (core_domain.User, error) {
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
			return core_domain.User{}, fmt.Errorf("%w: user not found", core_errors.ErrNotFound)
		}
		return core_domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	if phone.Valid {
		user.Phone = &phone.String
	}
	user.Role = core_domain.Role(roleValue)
	return user, nil
}
