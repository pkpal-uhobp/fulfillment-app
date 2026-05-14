package auth_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *AuthService) Logout(
	ctx context.Context,
	input LogoutInput,
) error {
	if strings.TrimSpace(input.RefreshToken) == "" {
		return fmt.Errorf(
			"%w: empty refresh token",
			core_errors.ErrInvalidArgument,
		)
	}

	claims, err := s.parseToken(input.RefreshToken)
	if err != nil {
		return fmt.Errorf(
			"%w: invalid refresh token",
			core_errors.ErrUnauthorized,
		)
	}

	if claims.TokenType != core_domain.TokenTypeRefresh.String() {
		return fmt.Errorf(
			"%w: invalid token type",
			core_errors.ErrUnauthorized,
		)
	}

	deviceID, err := uuid.Parse(claims.DeviceID)
	if err != nil {
		return fmt.Errorf(
			"%w: invalid device id",
			core_errors.ErrUnauthorized,
		)
	}

	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		return s.repo.RevokeActiveTokensByDevice(
			ctx,
			claims.UserID,
			deviceID,
			"logout",
		)
	})
}
