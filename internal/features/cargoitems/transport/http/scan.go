package cargoitems_transport_http

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
)

func (h *CargoItemsHTTPHandler) ScanCargoItem(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	cargoItem, err := h.cargoItemsService.ScanCargoItem(
		r.Context(),
		user.ID,
		user.Role,
		core_http_utils.QueryString(r, "qr_code"),
	)
	if err != nil {
		response.ErrorResponse(err, "scan cargo item")
		return
	}

	response.JSONResponse(CargoItemResponse{CargoItem: cargoItem}, http.StatusOK)
}
