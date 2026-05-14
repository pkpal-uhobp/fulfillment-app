package auth_transport_http

import (
	"context"
	"net/http"

	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

type AuthHTTPHandler struct {
	log            *core_logger.Logger
	authService    AuthService
	authMiddleware core_http_middleware.Middleware
}

type AuthService interface {
	Register(
		ctx context.Context,
		input auth_service.RegisterInput,
	) (auth_service.UserDTO, auth_service.TokenPair, error)

	Login(
		ctx context.Context,
		input auth_service.LoginInput,
	) (auth_service.UserDTO, auth_service.TokenPair, error)

	Refresh(
		ctx context.Context,
		input auth_service.RefreshInput,
	) (auth_service.TokenPair, error)

	Logout(
		ctx context.Context,
		input auth_service.LogoutInput,
	) error

	GetMe(
		ctx context.Context,
		userID int64,
	) (auth_service.UserDTO, error)

	VerifyAccessToken(
		ctx context.Context,
		accessToken string,
	) (auth_service.AuthClaims, error)
}

func NewAuthHTTPHandler(
	log *core_logger.Logger,
	authService AuthService,
	authMiddleware core_http_middleware.Middleware,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		log:            log,
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodPost,
			"/auth/register",
			h.Register,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/auth/login",
			h.Login,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/auth/refresh",
			h.Refresh,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/auth/logout",
			h.Logout,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/auth/me",
			h.Me,
			nil,
			h.authMiddleware,
		),
	}
}
