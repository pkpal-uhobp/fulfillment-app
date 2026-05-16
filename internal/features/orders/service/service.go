package orders_service

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type OrdersRepository interface {
	CreateOrder(
		ctx context.Context,
		order core_domain.Order,
		cargoPlaces []core_domain.OrderCargoPlace,
		pickup *core_domain.PickupRequest,
		changedBy int64,
	) (core_domain.OrderDetails, error)
	ListOrders(
		ctx context.Context,
		filter core_domain.OrderFilter,
	) ([]core_domain.OrderDetails, error)
	GetOrder(
		ctx context.Context,
		orderID int64,
	) (core_domain.OrderDetails, error)
	ListOrderStatusHistory(
		ctx context.Context,
		orderID int64,
	) ([]core_domain.OrderStatusHistory, error)
	CancelOrder(
		ctx context.Context,
		orderID int64,
		changedBy int64,
		comment *string,
	) error
	UpdateOrderStatus(
		ctx context.Context,
		orderID int64,
		status string,
		changedBy int64,
		comment *string,
	) error
}

type PickupCalendar interface {
	EnsureDateAvailable(
		ctx context.Context,
		warehouseID int64,
		date string,
	) error
}

type OrdersService struct {
	repo           OrdersRepository
	pickupCalendar PickupCalendar
}

func NewOrdersService(repo OrdersRepository) *OrdersService {
	return &OrdersService{
		repo: repo,
	}
}

func NewOrdersServiceWithPickupCalendar(
	repo OrdersRepository,
	pickupCalendar PickupCalendar,
) *OrdersService {
	return &OrdersService{
		repo:           repo,
		pickupCalendar: pickupCalendar,
	}
}
