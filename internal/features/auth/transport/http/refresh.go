package auth_transport_http

import (
	"net/http"

	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

func (h *AuthHTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request RefreshRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid refresh request")
		return
	}

	tokens, err := h.authService.Refresh(r.Context(), auth_service.RefreshInput{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		response.ErrorResponse(err, "refresh token")
		return
	}

	response.JSONResponse(RefreshResponse{
		Tokens: tokens,
	}, http.StatusOK)
}
