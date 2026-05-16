package orders_service_tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
)

type calendarCall struct {
	WarehouseID int64
	Date        string
}

type fakePickupCalendarForOrders struct {
	calls []calendarCall
	err   error
}

func (f *fakePickupCalendarForOrders) EnsureDateAvailable(
	ctx context.Context,
	warehouseID int64,
	date string,
) error {
	f.calls = append(f.calls, calendarCall{
		WarehouseID: warehouseID,
		Date:        date,
	})
	return f.err
}

type fakeOrdersRepositoryForCalendar struct {
	createOrderCalled bool
	createdOrder      core_domain.Order
	createdPickup     *core_domain.PickupRequest
	createdPlaces     []core_domain.OrderCargoPlace
}

func (f *fakeOrdersRepositoryForCalendar) CreateOrder(
	ctx context.Context,
	order core_domain.Order,
	cargoPlaces []core_domain.OrderCargoPlace,
	pickup *core_domain.PickupRequest,
	changedBy int64,
) (core_domain.OrderDetails, error) {
	f.createOrderCalled = true
	f.createdOrder = order
	f.createdPickup = pickup
	f.createdPlaces = cargoPlaces

	now := time.Now()
	order.ID = 100
	order.CreatedAt = now
	order.UpdatedAt = now
	for i := range cargoPlaces {
		cargoPlaces[i].ID = int64(i + 1)
		cargoPlaces[i].OrderID = order.ID
		cargoPlaces[i].CreatedAt = now
	}
	if pickup != nil {
		pickup.ID = 1
		pickup.OrderID = order.ID
		pickup.Status = "created"
		pickup.CreatedAt = now
		pickup.UpdatedAt = now
	}

	return core_domain.OrderDetails{
		Order:       order,
		CargoPlaces: cargoPlaces,
		Pickup:      pickup,
	}, nil
}

func (f *fakeOrdersRepositoryForCalendar) ListOrders(
	ctx context.Context,
	filter core_domain.OrderFilter,
) ([]core_domain.OrderDetails, error) {
	return nil, fmt.Errorf("unexpected ListOrders call")
}

func (f *fakeOrdersRepositoryForCalendar) GetOrder(
	ctx context.Context,
	orderID int64,
) (core_domain.OrderDetails, error) {
	return core_domain.OrderDetails{}, fmt.Errorf("unexpected GetOrder call")
}

func (f *fakeOrdersRepositoryForCalendar) ListOrderStatusHistory(
	ctx context.Context,
	orderID int64,
) ([]core_domain.OrderStatusHistory, error) {
	return nil, fmt.Errorf("unexpected ListOrderStatusHistory call")
}

func (f *fakeOrdersRepositoryForCalendar) CancelOrder(
	ctx context.Context,
	orderID int64,
	changedBy int64,
	comment *string,
) error {
	return fmt.Errorf("unexpected CancelOrder call")
}

func (f *fakeOrdersRepositoryForCalendar) UpdateOrderStatus(
	ctx context.Context,
	orderID int64,
	status string,
	changedBy int64,
	comment *string,
) error {
	return fmt.Errorf("unexpected UpdateOrderStatus call")
}

func TestCreateOrderSelfDeliveryChecksPickupCalendar(t *testing.T) {
	repo := &fakeOrdersRepositoryForCalendar{}
	calendar := &fakePickupCalendarForOrders{}
	service := orders_service.NewOrdersServiceWithPickupCalendar(repo, calendar)

	selfDeliveryDate := "2026-06-20"
	_, err := service.CreateOrder(
		context.Background(),
		10,
		core_domain.RoleClient.String(),
		orders_service.CreateOrderInput{
			ReceivingWarehouseID:    11,
			DestinationWarehouseID: 22,
			ProductTypeID:          33,
			HandoverType:           core_domain.HandoverTypeSelfDelivery.String(),
			SelfDeliveryDate:       &selfDeliveryDate,
			CargoPlaces: []orders_service.CargoPlaceInput{
				{CargoPlaceTypeID: 1, Quantity: 2},
			},
		},
	)
	if err != nil {
		t.Fatalf("CreateOrder() error = %v, want nil", err)
	}

	if len(calendar.calls) != 1 {
		t.Fatalf("calendar calls = %d, want 1", len(calendar.calls))
	}
	if calendar.calls[0].WarehouseID != 11 {
		t.Fatalf("calendar warehouse_id = %d, want 11", calendar.calls[0].WarehouseID)
	}
	if calendar.calls[0].Date != selfDeliveryDate {
		t.Fatalf("calendar date = %q, want %q", calendar.calls[0].Date, selfDeliveryDate)
	}
	if !repo.createOrderCalled {
		t.Fatal("repo.CreateOrder() was not called")
	}
}

func TestCreateOrderSelfDeliveryRejectsClosedPickupCalendarDate(t *testing.T) {
	repo := &fakeOrdersRepositoryForCalendar{}
	calendar := &fakePickupCalendarForOrders{
		err: fmt.Errorf("%w: pickup calendar date is closed", core_errors.ErrConflict),
	}
	service := orders_service.NewOrdersServiceWithPickupCalendar(repo, calendar)

	selfDeliveryDate := "2026-06-20"
	_, err := service.CreateOrder(
		context.Background(),
		10,
		core_domain.RoleClient.String(),
		orders_service.CreateOrderInput{
			ReceivingWarehouseID:    11,
			DestinationWarehouseID: 22,
			ProductTypeID:          33,
			HandoverType:           core_domain.HandoverTypeSelfDelivery.String(),
			SelfDeliveryDate:       &selfDeliveryDate,
			CargoPlaces: []orders_service.CargoPlaceInput{
				{CargoPlaceTypeID: 1, Quantity: 2},
			},
		},
	)
	if err == nil {
		t.Fatal("CreateOrder() error = nil, want conflict")
	}
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("CreateOrder() error = %v, want ErrConflict", err)
	}
	if repo.createOrderCalled {
		t.Fatal("repo.CreateOrder() was called for closed date")
	}
	if len(calendar.calls) != 1 {
		t.Fatalf("calendar calls = %d, want 1", len(calendar.calls))
	}
}

func TestCreateOrderPickupChecksPickupCalendar(t *testing.T) {
	repo := &fakeOrdersRepositoryForCalendar{}
	calendar := &fakePickupCalendarForOrders{}
	service := orders_service.NewOrdersServiceWithPickupCalendar(repo, calendar)

	_, err := service.CreateOrder(
		context.Background(),
		10,
		core_domain.RoleClient.String(),
		orders_service.CreateOrderInput{
			ReceivingWarehouseID:    11,
			DestinationWarehouseID: 22,
			ProductTypeID:          33,
			HandoverType:           core_domain.HandoverTypePickup.String(),
			Pickup: &orders_service.PickupInput{
				PickupAddress: "Москва, склад клиента",
				PickupDate:    "2026-06-21",
			},
			CargoPlaces: []orders_service.CargoPlaceInput{
				{CargoPlaceTypeID: 1, Quantity: 2},
			},
		},
	)
	if err != nil {
		t.Fatalf("CreateOrder() error = %v, want nil", err)
	}

	if len(calendar.calls) != 1 {
		t.Fatalf("calendar calls = %d, want 1", len(calendar.calls))
	}
	if calendar.calls[0].WarehouseID != 11 {
		t.Fatalf("calendar warehouse_id = %d, want 11", calendar.calls[0].WarehouseID)
	}
	if calendar.calls[0].Date != "2026-06-21" {
		t.Fatalf("calendar date = %q, want %q", calendar.calls[0].Date, "2026-06-21")
	}
	if !repo.createOrderCalled {
		t.Fatal("repo.CreateOrder() was not called")
	}
	if repo.createdPickup == nil {
		t.Fatal("pickup request was not passed to repository")
	}
}

func TestCreateOrderPickupRejectsClosedPickupCalendarDate(t *testing.T) {
	repo := &fakeOrdersRepositoryForCalendar{}
	calendar := &fakePickupCalendarForOrders{
		err: fmt.Errorf("%w: pickup calendar date is closed", core_errors.ErrConflict),
	}
	service := orders_service.NewOrdersServiceWithPickupCalendar(repo, calendar)

	_, err := service.CreateOrder(
		context.Background(),
		10,
		core_domain.RoleClient.String(),
		orders_service.CreateOrderInput{
			ReceivingWarehouseID:    11,
			DestinationWarehouseID: 22,
			ProductTypeID:          33,
			HandoverType:           core_domain.HandoverTypePickup.String(),
			Pickup: &orders_service.PickupInput{
				PickupAddress: "Москва, склад клиента",
				PickupDate:    "2026-06-21",
			},
			CargoPlaces: []orders_service.CargoPlaceInput{
				{CargoPlaceTypeID: 1, Quantity: 2},
			},
		},
	)
	if err == nil {
		t.Fatal("CreateOrder() error = nil, want conflict")
	}
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("CreateOrder() error = %v, want ErrConflict", err)
	}
	if repo.createOrderCalled {
		t.Fatal("repo.CreateOrder() was called for closed date")
	}
	if len(calendar.calls) != 1 {
		t.Fatalf("calendar calls = %d, want 1", len(calendar.calls))
	}
}
