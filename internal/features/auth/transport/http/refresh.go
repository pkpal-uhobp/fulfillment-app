package auth_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)


func (h *AuthHTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.ErrorResponse(
			fmt.Errorf("%w: decode refresh request", core_errors.ErrInvalidArgument),
			"invalid request body",
		)
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
