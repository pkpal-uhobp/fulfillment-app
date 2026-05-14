package auth_transport_http

import (
	"net/http"

	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid login request")
		return
	}

	user, tokens, err := h.authService.Login(r.Context(), auth_service.LoginInput{
		Email:    request.Email,
		Password: request.Password,
		DeviceID: request.DeviceID,
	})
	if err != nil {
		response.ErrorResponse(err, "login user")
		return
	}

	response.JSONResponse(LoginResponse{
		User:   user,
		Tokens: tokens,
	}, http.StatusOK)
}
