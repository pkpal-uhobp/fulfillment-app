package auth_service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	input RefreshInput,
) (TokenPair, error) {
	if strings.TrimSpace(input.RefreshToken) == "" {
		return TokenPair{}, fmt.Errorf(
			"%w: empty refresh token",
			core_errors.ErrInvalidArgument,
		)
	}

	claims, err := s.parseToken(input.RefreshToken)
	if err != nil {
		return TokenPair{}, fmt.Errorf(
			"%w: invalid refresh token",
			core_errors.ErrUnauthorized,
		)
	}

	if claims.TokenType != core_domain.TokenTypeRefresh.String() {
		return TokenPair{}, fmt.Errorf(
			"%w: invalid token type",
			core_errors.ErrUnauthorized,
		)
	}

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return TokenPair{}, fmt.Errorf(
			"%w: invalid token jti",
			core_errors.ErrUnauthorized,
		)
	}

	deviceID, err := uuid.Parse(claims.DeviceID)
	if err != nil {
		return TokenPair{}, fmt.Errorf(
			"%w: invalid device id",
			core_errors.ErrUnauthorized,
		)
	}

	var tokens TokenPair
	err = s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		issuedToken, err := s.repo.GetIssuedTokenByJTI(ctx, jti)
		if err != nil {
			if errors.Is(err, core_errors.ErrNotFound) {
				return fmt.Errorf(
					"%w: refresh token not found",
					core_errors.ErrUnauthorized,
				)
			}

			return err
		}

		if issuedToken.UserID != claims.UserID {
			return fmt.Errorf(
				"%w: invalid refresh token owner",
				core_errors.ErrUnauthorized,
			)
		}

		if issuedToken.TokenType != core_domain.TokenTypeRefresh ||
			issuedToken.Revoked ||
			time.Now().UTC().After(issuedToken.ExpiresAt) {
			return fmt.Errorf(
				"%w: refresh token is not active",
				core_errors.ErrUnauthorized,
			)
		}

		user, err := s.repo.GetUserByID(ctx, issuedToken.UserID)
		if err != nil {
			return err
		}

		if err := checkUserAvailable(user); err != nil {
			return err
		}

		if err := s.repo.RevokeTokenByJTI(ctx, jti, "refresh token rotated"); err != nil {
			if errors.Is(err, core_errors.ErrNotFound) {
				return fmt.Errorf(
					"%w: refresh token is not active",
					core_errors.ErrUnauthorized,
				)
			}

			return err
		}

		tokenPair, err := s.issueTokenPair(ctx, user, deviceID)
		if err != nil {
			return err
		}

		tokens = tokenPair

		return nil
	})
	if err != nil {
		return TokenPair{}, err
	}

	return tokens, nil
}
