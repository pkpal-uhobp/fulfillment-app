package auth_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
)

func (t *Transport) meRoute() core_http_server.Route {
	return core_http_server.NewRoute(
		http.MethodGet,
		"/auth/me",
		t.Me,
		nil,
		t.authMiddleware,
	)
}

func (t *Transport) Me(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(t.log, w)

	currentUser, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(
			fmt.Errorf("%w: current user not found", core_errors.ErrUnauthorized),
			"unauthorized",
		)
		return
	}

	user, err := t.service.GetMe(r.Context(), currentUser.ID)
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	response.JSONResponse(MeResponse{
		User: user,
	}, http.StatusOK)
}
