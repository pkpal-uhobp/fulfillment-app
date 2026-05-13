package auth_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *Repository) CreateIssuedToken(
	ctx context.Context,
	token core_domain.IssuedToken,
) error {
	q := r.tx.Querier(ctx)

	const query = `
		INSERT INTO issued_tokens (
			user_id,
			jti,
			token_type,
			device_id,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := q.Exec(
		ctx,
		query,
		token.UserID,
		token.JTI,
		token.TokenType.String(),
		token.DeviceID,
		token.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("create issued token: %w", err)
	}

	return nil
}
