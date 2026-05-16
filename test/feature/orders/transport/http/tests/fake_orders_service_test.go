package orders_http_tests

import (
	"context"

	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
)

type fakeOrdersService struct {
	createOrderFn       func(context.Context, int64, string, orders_service.CreateOrderInput) (orders_service.OrderDTO, error)
	listOrdersFn        func(context.Context, int64, string, orders_service.OrderFilter) ([]orders_service.OrderDTO, error)
	getOrderFn          func(context.Context, int64, int64, string) (orders_service.OrderDTO, error)
	getOrderHistoryFn   func(context.Context, int64, int64, string) ([]orders_service.OrderStatusHistoryDTO, error)
	cancelOrderFn       func(context.Context, int64, int64, string, orders_service.CancelOrderInput) error
	updateOrderStatusFn func(context.Context, int64, int64, orders_service.UpdateOrderStatusInput) (orders_service.OrderDTO, error)
}

func (f *fakeOrdersService) CreateOrder(
	ctx context.Context,
	actorID int64,
	actorRole string,
	input orders_service.CreateOrderInput,
) (orders_service.OrderDTO, error) {
	if f.createOrderFn == nil {
		return orders_service.OrderDTO{}, nil
	}
	return f.createOrderFn(ctx, actorID, actorRole, input)
}

func (f *fakeOrdersService) ListOrders(
	ctx context.Context,
	actorID int64,
	actorRole string,
	filter orders_service.OrderFilter,
) ([]orders_service.OrderDTO, error) {
	if f.listOrdersFn == nil {
		return nil, nil
	}
	return f.listOrdersFn(ctx, actorID, actorRole, filter)
}

func (f *fakeOrdersService) GetOrder(
	ctx context.Context,
	orderID int64,
	actorID int64,
	actorRole string,
) (orders_service.OrderDTO, error) {
	if f.getOrderFn == nil {
		return orders_service.OrderDTO{}, nil
	}
	return f.getOrderFn(ctx, orderID, actorID, actorRole)
}

func (f *fakeOrdersService) GetOrderHistory(
	ctx context.Context,
	orderID int64,
	actorID int64,
	actorRole string,
) ([]orders_service.OrderStatusHistoryDTO, error) {
	if f.getOrderHistoryFn == nil {
		return nil, nil
	}
	return f.getOrderHistoryFn(ctx, orderID, actorID, actorRole)
}

func (f *fakeOrdersService) CancelOrder(
	ctx context.Context,
	orderID int64,
	actorID int64,
	actorRole string,
	input orders_service.CancelOrderInput,
) error {
	if f.cancelOrderFn == nil {
		return nil
	}
	return f.cancelOrderFn(ctx, orderID, actorID, actorRole, input)
}

func (f *fakeOrdersService) UpdateOrderStatus(
	ctx context.Context,
	orderID int64,
	actorID int64,
	input orders_service.UpdateOrderStatusInput,
) (orders_service.OrderDTO, error) {
	if f.updateOrderStatusFn == nil {
		return orders_service.OrderDTO{}, nil
	}
	return f.updateOrderStatusFn(ctx, orderID, actorID, input)
}
