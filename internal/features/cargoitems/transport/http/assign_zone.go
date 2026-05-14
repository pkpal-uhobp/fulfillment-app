package cargoitems_transport_http

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
	cargoitems_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/service"
)

func (h *CargoItemsHTTPHandler) AssignStorageZone(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	cargoItemID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid cargo item id")
		return
	}

	var request AssignStorageZoneRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid assign storage zone request")
		return
	}

	cargoItem, err := h.cargoItemsService.AssignStorageZone(
		r.Context(),
		cargoItemID,
		user.ID,
		user.Role,
		cargoitems_service.AssignStorageZoneInput{
			StorageZoneID: request.StorageZoneID,
			Comment:       request.Comment,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "assign storage zone")
		return
	}

	response.JSONResponse(CargoItemResponse{CargoItem: cargoItem}, http.StatusOK)
}
