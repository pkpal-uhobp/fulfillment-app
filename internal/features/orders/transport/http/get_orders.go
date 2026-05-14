package orders_transport_http

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
)

func (h *OrdersHTTPHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	clientID, err := core_http_utils.QueryInt64Ptr(r, "client_id")
	if err != nil {
		response.ErrorResponse(err, "invalid client id")
		return
	}

	warehouseID, err := core_http_utils.QueryInt64Ptr(r, "warehouse_id")
	if err != nil {
		response.ErrorResponse(err, "invalid warehouse id")
		return
	}

	receivingWarehouseID, err := core_http_utils.QueryInt64Ptr(r, "receiving_warehouse_id")
	if err != nil {
		response.ErrorResponse(err, "invalid receiving warehouse id")
		return
	}

	destinationWarehouseID, err := core_http_utils.QueryInt64Ptr(r, "destination_warehouse_id")
	if err != nil {
		response.ErrorResponse(err, "invalid destination warehouse id")
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

	orders, err := h.ordersService.ListOrders(
		r.Context(),
		user.ID,
		user.Role,
		orders_service.OrderFilter{
			ClientID:               clientID,
			Status:                 core_http_utils.QueryString(r, "status"),
			HandoverType:           core_http_utils.QueryString(r, "handover_type"),
			WarehouseID:            warehouseID,
			ReceivingWarehouseID:   receivingWarehouseID,
			DestinationWarehouseID: destinationWarehouseID,
			Page:                   page,
			Limit:                  limit,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "list orders")
		return
	}

	response.JSONResponse(
		OrdersResponse{
			Orders: orders,
		},
		http.StatusOK,
	)
}

func queryInt(r *http.Request, name string) (int, error) {
	value, err := core_http_utils.QueryInt64Ptr(r, name)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return 0, nil
	}
	return int(*value), nil
}
