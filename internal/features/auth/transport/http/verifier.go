package auth_transport_http

import (
	"context"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
)

type AccessTokenVerifier struct {
	authService AuthService
}

func NewAccessTokenVerifier(authService AuthService) *AccessTokenVerifier {
	return &AccessTokenVerifier{
		authService: authService,
	}
}

func (v *AccessTokenVerifier) VerifyAccessToken(
	ctx context.Context,
	token string,
) (*core_http_middleware.AccessTokenClaims, error) {
	claims, err := v.authService.VerifyAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &core_http_middleware.AccessTokenClaims{
		UserID: claims.UserID,
		Role:   claims.Role,
		JTI:    claims.JTI,
	}, nil
}
