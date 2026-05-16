package orders_service_tests

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type fakeOrdersRepository struct {
	createOrderFn            func(context.Context, core_domain.Order, []core_domain.OrderCargoPlace, *core_domain.PickupRequest, int64) (core_domain.OrderDetails, error)
	listOrdersFn             func(context.Context, core_domain.OrderFilter) ([]core_domain.OrderDetails, error)
	getOrderFn               func(context.Context, int64) (core_domain.OrderDetails, error)
	listOrderStatusHistoryFn func(context.Context, int64) ([]core_domain.OrderStatusHistory, error)
	cancelOrderFn            func(context.Context, int64, int64, *string) error
	updateOrderStatusFn      func(context.Context, int64, string, int64, *string) error
}

func (f *fakeOrdersRepository) CreateOrder(ctx context.Context, order core_domain.Order, cargoPlaces []core_domain.OrderCargoPlace, pickup *core_domain.PickupRequest, changedBy int64) (core_domain.OrderDetails, error) {
	if f.createOrderFn == nil {
		return core_domain.OrderDetails{}, fmt.Errorf("unexpected CreateOrder call")
	}
	return f.createOrderFn(ctx, order, cargoPlaces, pickup, changedBy)
}

func (f *fakeOrdersRepository) ListOrders(ctx context.Context, filter core_domain.OrderFilter) ([]core_domain.OrderDetails, error) {
	if f.listOrdersFn == nil {
		return nil, fmt.Errorf("unexpected ListOrders call")
	}
	return f.listOrdersFn(ctx, filter)
}

func (f *fakeOrdersRepository) GetOrder(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
	if f.getOrderFn == nil {
		return core_domain.OrderDetails{}, fmt.Errorf("unexpected GetOrder call")
	}
	return f.getOrderFn(ctx, orderID)
}

func (f *fakeOrdersRepository) ListOrderStatusHistory(ctx context.Context, orderID int64) ([]core_domain.OrderStatusHistory, error) {
	if f.listOrderStatusHistoryFn == nil {
		return nil, fmt.Errorf("unexpected ListOrderStatusHistory call")
	}
	return f.listOrderStatusHistoryFn(ctx, orderID)
}

func (f *fakeOrdersRepository) CancelOrder(ctx context.Context, orderID int64, changedBy int64, comment *string) error {
	if f.cancelOrderFn == nil {
		return fmt.Errorf("unexpected CancelOrder call")
	}
	return f.cancelOrderFn(ctx, orderID, changedBy, comment)
}

func (f *fakeOrdersRepository) UpdateOrderStatus(ctx context.Context, orderID int64, status string, changedBy int64, comment *string) error {
	if f.updateOrderStatusFn == nil {
		return fmt.Errorf("unexpected UpdateOrderStatus call")
	}
	return f.updateOrderStatusFn(ctx, orderID, status, changedBy, comment)
}

func ptr[T any](value T) *T {
	return &value
}

func makeOrderDetails(id int64, clientID int64, handoverType core_domain.HandoverType, status core_domain.OrderStatus) core_domain.OrderDetails {
	return core_domain.OrderDetails{
		Order: core_domain.Order{
			ID:                     id,
			ClientID:               clientID,
			ReceivingWarehouseID:   10,
			DestinationWarehouseID: 20,
			ProductTypeID:          30,
			HandoverType:           handoverType,
			Status:                 status,
		},
		CargoPlaces: []core_domain.OrderCargoPlace{
			{ID: 1, OrderID: id, CargoPlaceTypeID: 7, Quantity: 2},
		},
	}
}
