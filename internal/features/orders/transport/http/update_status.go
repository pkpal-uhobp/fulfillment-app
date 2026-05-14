package orders_transport_http

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
)

func (h *OrdersHTTPHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	orderID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid order id")
		return
	}

	var request UpdateOrderStatusRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid update order status request")
		return
	}

	order, err := h.ordersService.UpdateOrderStatus(
		r.Context(),
		orderID,
		user.ID,
		orders_service.UpdateOrderStatusInput{
			Status:  request.Status,
			Comment: request.Comment,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "update order status")
		return
	}

	response.JSONResponse(
		OrderResponse{
			Order: order,
		},
		http.StatusOK,
	)
}
