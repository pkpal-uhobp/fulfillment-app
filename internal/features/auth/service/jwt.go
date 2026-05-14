package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type tokenClaims struct {
	UserID    int64  `json:"user_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	DeviceID  string `json:"device_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) generateToken(
	userID int64,
	role core_domain.Role,
	tokenType core_domain.TokenType,
	deviceID uuid.UUID,
	ttl time.Duration,
) (string, uuid.UUID, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	jti := uuid.New()

	claims := tokenClaims{
		UserID:    userID,
		Role:      role.String(),
		TokenType: tokenType.String(),
		DeviceID:  deviceID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti.String(),
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", uuid.Nil, time.Time{}, fmt.Errorf("sign JWT: %w", err)
	}

	return signedToken, jti, expiresAt, nil
}

func (s *AuthService) parseToken(tokenString string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&tokenClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(s.config.JWTSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("parse JWT: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid JWT claims")
	}

	return claims, nil
}

func (s *AuthService) issueTokenPair(
	ctx context.Context,
	user core_domain.User,
	deviceID uuid.UUID,
) (TokenPair, error) {
	accessToken, accessJTI, accessExpiresAt, err := s.generateToken(
		user.ID,
		user.Role,
		core_domain.TokenTypeAccess,
		deviceID,
		s.config.AccessTTL,
	)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, refreshJTI, refreshExpiresAt, err := s.generateToken(
		user.ID,
		user.Role,
		core_domain.TokenTypeRefresh,
		deviceID,
		s.config.RefreshTTL,
	)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.repo.CreateIssuedToken(ctx, core_domain.IssuedToken{
		UserID:    user.ID,
		JTI:       accessJTI,
		TokenType: core_domain.TokenTypeAccess,
		DeviceID:  deviceID,
		ExpiresAt: accessExpiresAt,
	}); err != nil {
		return TokenPair{}, err
	}

	if err := s.repo.CreateIssuedToken(ctx, core_domain.IssuedToken{
		UserID:    user.ID,
		JTI:       refreshJTI,
		TokenType: core_domain.TokenTypeRefresh,
		DeviceID:  deviceID,
		ExpiresAt: refreshExpiresAt,
	}); err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		DeviceID:         deviceID.String(),
	}, nil
}
