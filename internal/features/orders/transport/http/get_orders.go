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

	orders, err := h.ordersService.ListOrders(
		r.Context(),
		user.ID,
		user.Role,
		orders_service.OrderFilter{
			ClientID:     clientID,
			Status:      core_http_utils.QueryString(r, "status"),
			HandoverType: core_http_utils.QueryString(r, "handover_type"),
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
