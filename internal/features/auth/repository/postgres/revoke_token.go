package auth_repository_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *AuthRepository) RevokeTokenByJTI(
	ctx context.Context,
	jti uuid.UUID,
	reason string,
) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		UPDATE issued_tokens
		SET revoked = TRUE,
			revoked_at = CURRENT_TIMESTAMP,
			revoked_reason = $2
		WHERE jti = $1
			AND revoked = FALSE;
	`

	tag, err := q.Exec(ctx, query, jti, reason)
	if err != nil {
		return fmt.Errorf("revoke token by jti: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: token not found or already revoked", core_errors.ErrNotFound)
	}

	return nil
}

func (r *AuthRepository) RevokeActiveTokensByDevice(
	ctx context.Context,
	userID int64,
	deviceID uuid.UUID,
	reason string,
) error {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		UPDATE issued_tokens
		SET revoked = TRUE,
			revoked_at = CURRENT_TIMESTAMP,
			revoked_reason = $3
		WHERE user_id = $1
			AND device_id = $2
			AND revoked = FALSE;
	`

	_, err := q.Exec(ctx, query, userID, deviceID, reason)
	if err != nil {
		return fmt.Errorf("revoke active tokens by device: %w", err)
	}

	return nil
}
