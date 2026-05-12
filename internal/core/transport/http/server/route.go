package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.Handler
	Roles       []string
	Middlewares []core_http_middleware.Middleware
}

func NewRoute(
	method string,
	path string,
	handler http.HandlerFunc,
	roles []string,
	middlewares ...core_http_middleware.Middleware,
) Route {
	return Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Roles:       roles,
		Middlewares: middlewares,
	}
}

func NewHandlerRoute(
	method string,
	path string,
	handler http.Handler,
	roles []string,
	middlewares ...core_http_middleware.Middleware,
) Route {
	return Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Roles:       roles,
		Middlewares: middlewares,
	}
}
