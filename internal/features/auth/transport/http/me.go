package auth_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
)

func (h *AuthHTTPHandler) Me(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	currentUser, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(
			fmt.Errorf("%w: current user not found", core_errors.ErrUnauthorized),
			"unauthorized",
		)
		return
	}

	user, err := h.authService.GetMe(r.Context(), currentUser.ID)
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	response.JSONResponse(MeResponse{
		User: user,
	}, http.StatusOK)
}
