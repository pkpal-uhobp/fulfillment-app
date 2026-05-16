package pickupcalendar_transport_http

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
	pickupcalendar_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/service"
)

func (h *PickupCalendarHTTPHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	warehouseID, err := core_http_utils.QueryInt64(r, "warehouse_id")
	if err != nil {
		response.ErrorResponse(err, "invalid warehouse id")
		return
	}
	days, err := h.pickupCalendarService.GetCalendar(
		r.Context(),
		user.ID,
		user.Role,
		pickupcalendar_service.CalendarFilter{
			WarehouseID: warehouseID,
			DateFrom:    core_http_utils.QueryString(r, "date_from"),
			DateTo:      core_http_utils.QueryString(r, "date_to"),
		},
	)
	if err != nil {
		response.ErrorResponse(err, "get pickup calendar")
		return
	}
	response.JSONResponse(PickupCalendarResponse{Days: days}, http.StatusOK)
}

func (h *PickupCalendarHTTPHandler) BlockDate(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	var request BlockDateRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid block pickup date request")
		return
	}
	block, err := h.pickupCalendarService.BlockDate(
		r.Context(),
		user.ID,
		user.Role,
		pickupcalendar_service.BlockDateInput{
			WarehouseID: request.WarehouseID,
			BlockedDate: request.BlockedDate,
			Reason:      request.Reason,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "block pickup date")
		return
	}
	response.JSONResponse(PickupCalendarBlockResponse{Block: block}, http.StatusCreated)
}

func (h *PickupCalendarHTTPHandler) UnblockDate(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	blockID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid pickup calendar block id")
		return
	}
	if err := h.pickupCalendarService.UnblockDate(r.Context(), user.ID, user.Role, blockID); err != nil {
		response.ErrorResponse(err, "unblock pickup date")
		return
	}
	response.NoContentResponse()
}

func (h *PickupCalendarHTTPHandler) SetCapacity(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	var request SetCapacityRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid set pickup capacity request")
		return
	}
	capacity, err := h.pickupCalendarService.SetCapacity(
		r.Context(),
		user.ID,
		user.Role,
		pickupcalendar_service.SetCapacityInput{
			WarehouseID:   request.WarehouseID,
			PickupDate:    request.PickupDate,
			MaxOrders:     request.MaxOrders,
			CurrentOrders: request.CurrentOrders,
			IsClosed:      request.IsClosed,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "set pickup capacity")
		return
	}
	response.JSONResponse(PickupCapacityResponse{Capacity: capacity}, http.StatusOK)
}
