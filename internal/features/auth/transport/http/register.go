package auth_transport_http

import (
	"net/http"

	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid register request")
		return
	}

	user, tokens, err := h.authService.Register(r.Context(), auth_service.RegisterInput{
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
