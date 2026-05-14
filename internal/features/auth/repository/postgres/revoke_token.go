package auth_repository_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *AuthRepository) RevokeTokenByJTI(
	ctx context.Context,
	jti uuid.UUID,
	reason string,
) error {
	q := r.tx.Querier(ctx)

	const query = `
		UPDATE issued_tokens
		SET
			revoked = TRUE,
			revoked_at = CURRENT_TIMESTAMP,
			revoked_reason = $2
		WHERE jti = $1
		  AND revoked = FALSE;
	`

	_, err := q.Exec(ctx, query, jti, reason)
	if err != nil {
		return fmt.Errorf("revoke token by jti: %w", err)
	}

	return nil
}

func (r *AuthRepository) RevokeActiveTokensByDevice(
	ctx context.Context,
	userID int64,
	deviceID uuid.UUID,
	reason string,
) error {
	q := r.tx.Querier(ctx)

	const query = `
		UPDATE issued_tokens
		SET
			revoked = TRUE,
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
