package auth_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *AuthRepository) GetIssuedTokenByJTI(
	ctx context.Context,
	jti uuid.UUID,
) (core_domain.IssuedToken, error) {
	q := r.tx.Querier(ctx)

	const query = `
		SELECT
			id,
			user_id,
			jti,
			token_type,
			device_id,
			revoked,
			expires_at
		FROM issued_tokens
		WHERE jti = $1;
	`

	var (
		token          core_domain.IssuedToken
		tokenTypeValue string
	)

	err := q.QueryRow(ctx, query, jti).Scan(
		&token.ID,
		&token.UserID,
		&token.JTI,
		&tokenTypeValue,
		&token.DeviceID,
		&token.Revoked,
		&token.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.IssuedToken{}, fmt.Errorf(
				"%w: token not found",
				core_errors.ErrNotFound,
			)
		}

		return core_domain.IssuedToken{}, fmt.Errorf("get issued token by jti: %w", err)
	}

	token.TokenType = core_domain.TokenType(tokenTypeValue)

	return token, nil
}
