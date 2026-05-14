package cargoitems_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
	cargoitems_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/service"
)

func (h *CargoItemsHTTPHandler) ListCargoItems(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	orderID, err := core_http_utils.QueryInt64Ptr(r, "order_id")
	if err != nil {
		response.ErrorResponse(err, "invalid order id")
		return
	}

	storageZoneID, err := core_http_utils.QueryInt64Ptr(r, "storage_zone_id")
	if err != nil {
		response.ErrorResponse(err, "invalid storage zone id")
		return
	}
	if storageZoneID == nil {
		storageZoneID, err = core_http_utils.QueryInt64Ptr(r, "zone_id")
		if err != nil {
			response.ErrorResponse(err, "invalid zone id")
			return
		}
	}

	gateID, err := core_http_utils.QueryInt64Ptr(r, "gate_id")
	if err != nil {
		response.ErrorResponse(err, "invalid gate id")
		return
	}

	page, err := queryInt(r, "page")
	if err != nil {
		response.ErrorResponse(err, "invalid page")
		return
	}
	limit, err := queryInt(r, "limit")
	if err != nil {
		response.ErrorResponse(err, "invalid limit")
		return
	}

	cargoItems, err := h.cargoItemsService.ListCargoItems(
		r.Context(),
		user.ID,
		user.Role,
		cargoitems_service.CargoItemFilter{
			OrderID:       orderID,
			Status:        core_http_utils.QueryString(r, "status"),
			StorageZoneID: storageZoneID,
			GateID:        gateID,
			QRCode:        core_http_utils.QueryString(r, "qr_code"),
			Page:          page,
			Limit:         limit,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "list cargo items")
		return
	}

	response.JSONResponse(CargoItemsResponse{CargoItems: cargoItems}, http.StatusOK)
}

func queryInt(r *http.Request, name string) (int, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%w: invalid query param %s", core_errors.ErrInvalidArgument, name)
	}

	return parsed, nil
}
