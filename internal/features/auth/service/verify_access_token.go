package auth_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *Service) VerifyAccessToken(
	ctx context.Context,
	accessToken string,
) (AuthClaims, error) {
	if strings.TrimSpace(accessToken) == "" {
		return AuthClaims{}, fmt.Errorf(
			"%w: empty access token",
			core_errors.ErrUnauthorized,
		)
	}

	claims, err := s.parseToken(accessToken)
	if err != nil {
		return AuthClaims{}, fmt.Errorf(
			"%w: invalid access token",
			core_errors.ErrUnauthorized,
		)
	}

	if claims.TokenType != core_domain.TokenTypeAccess.String() {
		return AuthClaims{}, fmt.Errorf(
			"%w: invalid token type",
			core_errors.ErrUnauthorized,
		)
	}

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return AuthClaims{}, fmt.Errorf(
			"%w: invalid token jti",
			core_errors.ErrUnauthorized,
		)
	}

	issuedToken, err := s.repo.GetIssuedTokenByJTI(ctx, jti)
	if err != nil {
		return AuthClaims{}, fmt.Errorf(
			"%w: access token not found",
			core_errors.ErrUnauthorized,
		)
	}

	if issuedToken.UserID != claims.UserID {
		return AuthClaims{}, fmt.Errorf(
			"%w: invalid access token owner",
			core_errors.ErrUnauthorized,
		)
	}

	if issuedToken.TokenType != core_domain.TokenTypeAccess ||
		issuedToken.Revoked ||
		time.Now().UTC().After(issuedToken.ExpiresAt) {
		return AuthClaims{}, fmt.Errorf(
			"%w: access token is not active",
			core_errors.ErrUnauthorized,
		)
	}

	user, err := s.repo.GetUserByID(ctx, issuedToken.UserID)
	if err != nil {
		return AuthClaims{}, err
	}

	if err := checkUserAvailable(user); err != nil {
		return AuthClaims{}, err
	}

	return AuthClaims{
		UserID:   user.ID,
		Role:     user.Role.String(),
		JTI:      claims.ID,
		DeviceID: claims.DeviceID,
	}, nil
}
