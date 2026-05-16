package users_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *UsersRepository) ListUsers(
	ctx context.Context,
	filter core_domain.UserFilter,
) ([]core_domain.User, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)
	offset := (filter.Page - 1) * filter.Limit
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE ($1 = '' OR role = $1)
		  AND ($2::boolean IS NULL OR is_active = $2)
		  AND ($3::boolean IS NULL OR is_blocked = $3)
		  AND (
			$4 = ''
			OR lower(email) LIKE '%%' || lower($4) || '%%'
			OR lower(full_name) LIKE '%%' || lower($4) || '%%'
		  )
		ORDER BY id
		LIMIT $5 OFFSET $6;
	`, userColumns)

	rows, err := q.Query(
		ctx,
		query,
		filter.Role,
		filter.IsActive,
		filter.IsBlocked,
		filter.Search,
		filter.Limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]core_domain.User, 0)
	for rows.Next() {
		user, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return users, nil
}
