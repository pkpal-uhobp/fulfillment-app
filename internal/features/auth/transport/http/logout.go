package auth_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

func (t *Transport) logoutRoute() core_http_server.Route {
	return core_http_server.NewRoute(
		http.MethodPost,
		"/auth/logout",
		t.Logout,
		nil,
	)
}

func (t *Transport) Logout(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(t.log, w)

	var request LogoutRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.ErrorResponse(
			fmt.Errorf("%w: decode logout request", core_errors.ErrInvalidArgument),
			"invalid request body",
		)
		return
	}

	if err := t.service.Logout(r.Context(), auth_service.LogoutInput{
		RefreshToken: request.RefreshToken,
	}); err != nil {
		response.ErrorResponse(err, "logout user")
		return
	}

	response.NoContentResponse()
}
