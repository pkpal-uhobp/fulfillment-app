package users_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	userID int64,
	patch core_domain.UserPatch,
) (core_domain.User, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	var (
		roleValue *string
	)
	if patch.Role != nil {
		value := patch.Role.String()
		roleValue = &value
	}

	query := fmt.Sprintf(`
		UPDATE users
		SET
			full_name = CASE WHEN $2 THEN $3 ELSE full_name END,
			phone = CASE WHEN $4 THEN $5 ELSE phone END,
			role = CASE WHEN $6 THEN $7 ELSE role END,
			is_active = CASE WHEN $8 THEN $9 ELSE is_active END,
			is_blocked = CASE WHEN $10 THEN $11 ELSE is_blocked END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING %s;
	`, userColumns)

	user, err := scanUserRow(q.QueryRow(
		ctx,
		query,
		userID,
		patch.FullName != nil,
		patch.FullName,
		patch.PhoneProvided,
		patch.Phone,
		patch.Role != nil,
		roleValue,
		patch.IsActive != nil,
		patch.IsActive,
		patch.IsBlocked != nil,
		patch.IsBlocked,
	))
	if err != nil {
		return core_domain.User{}, fmt.Errorf("patch user: %w", err)
	}
	return user, nil
}
