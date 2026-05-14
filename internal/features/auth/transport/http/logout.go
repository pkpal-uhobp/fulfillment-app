package auth_transport_http

import (
	"net/http"

	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request LogoutRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid logout request")
		return
	}

	if err := h.authService.Logout(r.Context(), auth_service.LogoutInput{
		RefreshToken: request.RefreshToken,
	}); err != nil {
		response.ErrorResponse(err, "logout user")
		return
	}

	response.NoContentResponse()
}
