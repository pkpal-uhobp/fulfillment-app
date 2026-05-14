package orders_transport_http

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
)

func (h *OrdersHTTPHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
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

	order, err := h.ordersService.GetOrder(
		r.Context(),
		orderID,
		user.ID,
		user.Role,
	)
	if err != nil {
		response.ErrorResponse(err, "get order")
		return
	}

	response.JSONResponse(
		OrderResponse{
			Order: order,
		},
		http.StatusOK,
	)
}
