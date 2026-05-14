package orders_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
)

type OrdersHTTPHandler struct {
	log           *core_logger.Logger
	ordersService OrdersService
}

type OrdersService interface {
	CreateOrder(
		ctx context.Context,
		actorID int64,
		actorRole string,
		input orders_service.CreateOrderInput,
	) (orders_service.OrderDTO, error)

	ListOrders(
		ctx context.Context,
		actorID int64,
		actorRole string,
		filter orders_service.OrderFilter,
	) ([]orders_service.OrderDTO, error)

	GetOrder(
		ctx context.Context,
		orderID int64,
		actorID int64,
		actorRole string,
	) (orders_service.OrderDTO, error)

	GetOrderHistory(
		ctx context.Context,
		orderID int64,
		actorID int64,
		actorRole string,
	) ([]orders_service.OrderStatusHistoryDTO, error)

	CancelOrder(
		ctx context.Context,
		orderID int64,
		actorID int64,
		actorRole string,
		input orders_service.CancelOrderInput,
	) error

	UpdateOrderStatus(
		ctx context.Context,
		orderID int64,
		actorID int64,
		input orders_service.UpdateOrderStatusInput,
	) (orders_service.OrderDTO, error)
}

func NewOrdersHTTPHandler(
	log *core_logger.Logger,
	ordersService OrdersService,
) *OrdersHTTPHandler {
	return &OrdersHTTPHandler{
		log:           log,
		ordersService: ordersService,
	}
}

func (h *OrdersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodPost,
			"/orders",
			h.CreateOrder,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/orders",
			h.ListOrders,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/orders/{id}",
			h.GetOrder,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/orders/{id}/history",
			h.GetOrderHistory,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/orders/{id}/cancel",
			h.CancelOrder,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/orders/{id}/status",
			h.UpdateOrderStatus,
			[]string{
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
	}
}
