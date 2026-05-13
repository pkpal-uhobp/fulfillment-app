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

func (t *Transport) registerRoute() core_http_server.Route {
	return core_http_server.NewRoute(
		http.MethodPost,
		"/auth/register",
		t.Register,
		nil,
	)
}

func (t *Transport) Register(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(t.log, w)

	var request RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.ErrorResponse(
			fmt.Errorf("%w: decode register request", core_errors.ErrInvalidArgument),
			"invalid request body",
		)
		return
	}

	user, tokens, err := t.service.Register(r.Context(), auth_service.RegisterInput{
		Email:    request.Email,
		Password: request.Password,
		FullName: request.FullName,
		Phone:    request.Phone,
	})
	if err != nil {
		response.ErrorResponse(err, "register user")
		return
	}

	response.JSONResponse(RegisterResponse{
		User:   user,
		Tokens: tokens,
	}, http.StatusCreated)
}
