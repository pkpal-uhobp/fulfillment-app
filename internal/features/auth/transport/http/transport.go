package auth_transport_http

import (
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

type Transport struct {
	log            *core_logger.Logger
	service        *auth_service.Service
	authMiddleware core_http_middleware.Middleware
}

func NewTransport(
	log *core_logger.Logger,
	service *auth_service.Service,
	authMiddleware core_http_middleware.Middleware,
) *Transport {
	return &Transport{
		log:            log,
		service:        service,
		authMiddleware: authMiddleware,
	}
}

func (t *Transport) RegisterRoutes(router *core_http_server.APIVersionRouter) {
	router.RegisterRoutes(
		t.registerRoute(),
		t.loginRoute(),
		t.refreshRoute(),
		t.logoutRoute(),
		t.meRoute(),
	)
}
