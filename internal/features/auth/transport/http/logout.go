package auth_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)


func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request LogoutRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.ErrorResponse(
			fmt.Errorf("%w: decode logout request", core_errors.ErrInvalidArgument),
			"invalid request body",
		)
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
