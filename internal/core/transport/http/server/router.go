package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type RoleMiddleware func(roles ...string) core_http_middleware.Middleware

type APIVersionRouter struct {
	*http.ServeMux

	apiVersion     ApiVersion
	middlewares    []core_http_middleware.Middleware
	roleMiddleware RoleMiddleware
}

func NewAPIVersionRouter(
	apiVersion ApiVersion,
	middlewares ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:    http.NewServeMux(),
		apiVersion:  apiVersion,
		middlewares: middlewares,
	}
}

func (r *APIVersionRouter) APIVersion() ApiVersion {
	return r.apiVersion
}

func (r *APIVersionRouter) SetRoleMiddleware(roleMiddleware RoleMiddleware) {
	r.roleMiddleware = roleMiddleware
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		handler := route.Handler

		routeMiddlewares := make(
			[]core_http_middleware.Middleware,
			0,
			len(route.Middlewares)+1,
		)
		
		routeMiddlewares = append(routeMiddlewares, route.Middlewares...)

		if len(route.Roles) > 0 {
			if r.roleMiddleware == nil {
				panic(fmt.Sprintf(
					"route %s %s has roles, but role middleware is not set",
					route.Method,
					route.Path,
				))
			}

			routeMiddlewares = append(
				routeMiddlewares,
				r.roleMiddleware(route.Roles...),
			)
		}

		if len(routeMiddlewares) > 0 {
			handler = core_http_middleware.ChainMiddlewares(
				handler,
				routeMiddlewares...,
			)
		}
		if len(r.middlewares) > 0 {
			handler = core_http_middleware.ChainMiddlewares(
				handler,
				r.middlewares...,
			)
		}

		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, handler)
	}
}

func (r *APIVersionRouter) RegisterRouters(routes ...Route) {
	r.RegisterRoutes(routes...)
}
