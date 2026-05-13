package auth_transport_http

import (
	"context"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

type AccessTokenVerifier struct {
	service *auth_service.Service
}

func NewAccessTokenVerifier(service *auth_service.Service) *AccessTokenVerifier {
	return &AccessTokenVerifier{
		service: service,
	}
}

func (v *AccessTokenVerifier) VerifyAccessToken(
	ctx context.Context,
	token string,
) (*core_http_middleware.AccessTokenClaims, error) {
	claims, err := v.service.VerifyAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &core_http_middleware.AccessTokenClaims{
		UserID: claims.UserID,
		Role:   claims.Role,
		JTI:    claims.JTI,
	}, nil
}
