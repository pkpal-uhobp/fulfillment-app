package orders_service_tests

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
)

func TestCreateOrderSelfDeliverySuccess(t *testing.T) {
	selfDeliveryDate := "2026-05-20"
	timeFrom := "10:00"
	timeTo := "12:00"
	comment := "Самопривоз клиентом"

	repo := &fakeOrdersRepository{
		createOrderFn: func(ctx context.Context, order core_domain.Order, cargoPlaces []core_domain.OrderCargoPlace, pickup *core_domain.PickupRequest, changedBy int64) (core_domain.OrderDetails, error) {
			if changedBy != 42 {
				t.Fatalf("changedBy = %d, want 42", changedBy)
			}
			if order.ClientID != 42 {
				t.Fatalf("ClientID = %d, want 42", order.ClientID)
			}
			if order.HandoverType != core_domain.HandoverTypeSelfDelivery {
				t.Fatalf("HandoverType = %q, want self_delivery", order.HandoverType)
			}
			if order.Status != core_domain.OrderStatusCreated {
				t.Fatalf("Status = %q, want created", order.Status)
			}
			if order.SelfDeliveryDate == nil || *order.SelfDeliveryDate != selfDeliveryDate {
				t.Fatalf("SelfDeliveryDate = %v, want %s", order.SelfDeliveryDate, selfDeliveryDate)
			}
			if order.SelfDeliveryTimeFrom == nil || *order.SelfDeliveryTimeFrom != timeFrom {
				t.Fatalf("SelfDeliveryTimeFrom = %v, want %s", order.SelfDeliveryTimeFrom, timeFrom)
			}
			if order.SelfDeliveryTimeTo == nil || *order.SelfDeliveryTimeTo != timeTo {
				t.Fatalf("SelfDeliveryTimeTo = %v, want %s", order.SelfDeliveryTimeTo, timeTo)
			}
			if order.Comment == nil || *order.Comment != comment {
				t.Fatalf("Comment = %v, want %s", order.Comment, comment)
			}
			if pickup != nil {
				t.Fatalf("pickup = %+v, want nil", pickup)
			}
			if len(cargoPlaces) != 1 || cargoPlaces[0].CargoPlaceTypeID != 7 || cargoPlaces[0].Quantity != 2 {
				t.Fatalf("unexpected cargoPlaces: %+v", cargoPlaces)
			}

			details := makeOrderDetails(100, order.ClientID, order.HandoverType, order.Status)
			details.Order.SelfDeliveryDate = order.SelfDeliveryDate
			details.Order.SelfDeliveryTimeFrom = order.SelfDeliveryTimeFrom
			details.Order.SelfDeliveryTimeTo = order.SelfDeliveryTimeTo
			details.Order.Comment = order.Comment
			return details, nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	dto, err := service.CreateOrder(context.Background(), 42, core_domain.RoleClient.String(), orders_service.CreateOrderInput{
		ReceivingWarehouseID:   10,
		DestinationWarehouseID: 20,
		ProductTypeID:          30,
		HandoverType:           core_domain.HandoverTypeSelfDelivery.String(),
		SelfDeliveryDate:       &selfDeliveryDate,
		SelfDeliveryTimeFrom:   &timeFrom,
		SelfDeliveryTimeTo:     &timeTo,
		Comment:                &comment,
		CargoPlaces: []orders_service.CargoPlaceInput{
			{CargoPlaceTypeID: 7, Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if dto.ID != 100 || dto.ClientID != 42 || dto.Status != core_domain.OrderStatusCreated.String() {
		t.Fatalf("unexpected dto: %+v", dto)
	}
	if len(dto.CargoPlaces) != 1 || dto.CargoPlaces[0].Quantity != 2 {
		t.Fatalf("unexpected cargo places dto: %+v", dto.CargoPlaces)
	}
}

func TestCreateOrderPickupSuccess(t *testing.T) {
	pickupDate := "2026-05-21"
	timeFrom := "09:00"
	timeTo := "11:00"
	contactName := "Иван Иванов"
	contactPhone := "+79990000000"

	repo := &fakeOrdersRepository{
		createOrderFn: func(ctx context.Context, order core_domain.Order, cargoPlaces []core_domain.OrderCargoPlace, pickup *core_domain.PickupRequest, changedBy int64) (core_domain.OrderDetails, error) {
			if order.HandoverType != core_domain.HandoverTypePickup {
				t.Fatalf("HandoverType = %q, want pickup", order.HandoverType)
			}
			if order.SelfDeliveryDate != nil || order.SelfDeliveryTimeFrom != nil || order.SelfDeliveryTimeTo != nil {
				t.Fatalf("self delivery fields must be nil: %+v", order)
			}
			if pickup == nil {
				t.Fatal("pickup = nil, want pickup request")
			}
			if pickup.PickupAddress != "Пермь, Ленина 1" || pickup.PickupDate != pickupDate {
				t.Fatalf("unexpected pickup: %+v", pickup)
			}
			if pickup.PickupTimeFrom == nil || *pickup.PickupTimeFrom != timeFrom || pickup.PickupTimeTo == nil || *pickup.PickupTimeTo != timeTo {
				t.Fatalf("unexpected pickup time range: %+v", pickup)
			}
			if pickup.ContactName == nil || *pickup.ContactName != contactName || pickup.ContactPhone == nil || *pickup.ContactPhone != contactPhone {
				t.Fatalf("unexpected contact data: %+v", pickup)
			}

			details := makeOrderDetails(101, order.ClientID, order.HandoverType, order.Status)
			details.Pickup = pickup
			details.Pickup.ID = 55
			details.Pickup.OrderID = 101
			details.Pickup.Status = "created"
			return details, nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	dto, err := service.CreateOrder(context.Background(), 42, core_domain.RoleClient.String(), orders_service.CreateOrderInput{
		ReceivingWarehouseID:   10,
		DestinationWarehouseID: 20,
		ProductTypeID:          30,
		HandoverType:           core_domain.HandoverTypePickup.String(),
		CargoPlaces: []orders_service.CargoPlaceInput{
			{CargoPlaceTypeID: 7, Quantity: 1},
		},
		Pickup: &orders_service.PickupInput{
			PickupAddress:  "Пермь, Ленина 1",
			PickupDate:     pickupDate,
			PickupTimeFrom: &timeFrom,
			PickupTimeTo:   &timeTo,
			ContactName:    &contactName,
			ContactPhone:   &contactPhone,
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if dto.Pickup == nil || dto.Pickup.PickupAddress != "Пермь, Ленина 1" {
		t.Fatalf("unexpected pickup dto: %+v", dto.Pickup)
	}
}

func TestCreateOrderRejectsForbiddenRole(t *testing.T) {
	service := orders_service.NewOrdersService(&fakeOrdersRepository{})
	_, err := service.CreateOrder(context.Background(), 5, core_domain.RoleWorker.String(), validSelfDeliveryInput())
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
}

func TestCreateOrderRejectsInvalidCargoPlaceQuantity(t *testing.T) {
	input := validSelfDeliveryInput()
	input.CargoPlaces[0].Quantity = 0

	service := orders_service.NewOrdersService(&fakeOrdersRepository{})
	_, err := service.CreateOrder(context.Background(), 42, core_domain.RoleClient.String(), input)
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestCreateOrderRejectsPickupWithSelfDeliverySchedule(t *testing.T) {
	selfDeliveryDate := "2026-05-20"
	pickupDate := "2026-05-21"
	input := orders_service.CreateOrderInput{
		ReceivingWarehouseID:   10,
		DestinationWarehouseID: 20,
		ProductTypeID:          30,
		HandoverType:           core_domain.HandoverTypePickup.String(),
		SelfDeliveryDate:       &selfDeliveryDate,
		CargoPlaces:            []orders_service.CargoPlaceInput{{CargoPlaceTypeID: 7, Quantity: 1}},
		Pickup: &orders_service.PickupInput{
			PickupAddress: "Пермь, Ленина 1",
			PickupDate:    pickupDate,
		},
	}

	service := orders_service.NewOrdersService(&fakeOrdersRepository{})
	_, err := service.CreateOrder(context.Background(), 42, core_domain.RoleClient.String(), input)
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestCreateOrderRejectsSelfDeliveryWithPickupData(t *testing.T) {
	input := validSelfDeliveryInput()
	input.Pickup = &orders_service.PickupInput{PickupAddress: "Пермь, Ленина 1", PickupDate: "2026-05-21"}

	service := orders_service.NewOrdersService(&fakeOrdersRepository{})
	_, err := service.CreateOrder(context.Background(), 42, core_domain.RoleClient.String(), input)
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestListOrdersScopesClientAndNormalizesPaging(t *testing.T) {
	wrongClientID := int64(999)
	receivingWarehouseID := int64(10)

	repo := &fakeOrdersRepository{
		listOrdersFn: func(ctx context.Context, filter core_domain.OrderFilter) ([]core_domain.OrderDetails, error) {
			if filter.ClientID == nil || *filter.ClientID != 42 {
				t.Fatalf("ClientID = %v, want 42", filter.ClientID)
			}
			if filter.Status != core_domain.OrderStatusStored.String() {
				t.Fatalf("Status = %q, want stored", filter.Status)
			}
			if filter.HandoverType != core_domain.HandoverTypeSelfDelivery.String() {
				t.Fatalf("HandoverType = %q, want self_delivery", filter.HandoverType)
			}
			if filter.ReceivingWarehouseID == nil || *filter.ReceivingWarehouseID != receivingWarehouseID {
				t.Fatalf("ReceivingWarehouseID = %v, want %d", filter.ReceivingWarehouseID, receivingWarehouseID)
			}
			if filter.Page != 1 || filter.Limit != 100 {
				t.Fatalf("Page/Limit = %d/%d, want 1/100", filter.Page, filter.Limit)
			}
			return []core_domain.OrderDetails{makeOrderDetails(100, 42, core_domain.HandoverTypeSelfDelivery, core_domain.OrderStatusStored)}, nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	orders, err := service.ListOrders(context.Background(), 42, core_domain.RoleClient.String(), orders_service.OrderFilter{
		ClientID:             &wrongClientID,
		Status:               core_domain.OrderStatusStored.String(),
		HandoverType:         core_domain.HandoverTypeSelfDelivery.String(),
		ReceivingWarehouseID: &receivingWarehouseID,
		Page:                 -5,
		Limit:                999,
	})
	if err != nil {
		t.Fatalf("ListOrders returned error: %v", err)
	}
	if len(orders) != 1 || orders[0].ID != 100 {
		t.Fatalf("unexpected orders: %+v", orders)
	}
}

func TestListOrdersRejectsInvalidStatus(t *testing.T) {
	service := orders_service.NewOrdersService(&fakeOrdersRepository{})
	_, err := service.ListOrders(context.Background(), 1, core_domain.RoleAdmin.String(), orders_service.OrderFilter{Status: "wrong"})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestGetOrderRejectsForeignClient(t *testing.T) {
	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			if orderID != 100 {
				t.Fatalf("orderID = %d, want 100", orderID)
			}
			return makeOrderDetails(orderID, 99, core_domain.HandoverTypePickup, core_domain.OrderStatusCreated), nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	_, err := service.GetOrder(context.Background(), 100, 42, core_domain.RoleClient.String())
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
}

func TestGetOrderAllowsAdminForAnyClient(t *testing.T) {
	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			return makeOrderDetails(orderID, 99, core_domain.HandoverTypePickup, core_domain.OrderStatusCreated), nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	dto, err := service.GetOrder(context.Background(), 100, 42, core_domain.RoleAdmin.String())
	if err != nil {
		t.Fatalf("GetOrder returned error: %v", err)
	}
	if dto.ClientID != 99 {
		t.Fatalf("ClientID = %d, want 99", dto.ClientID)
	}
}

func TestGetOrderHistoryChecksClientAccessAndMapsDTO(t *testing.T) {
	oldStatus := core_domain.OrderStatusCreated
	now := time.Now().UTC()

	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			return makeOrderDetails(orderID, 42, core_domain.HandoverTypePickup, core_domain.OrderStatusWaitingPickup), nil
		},
		listOrderStatusHistoryFn: func(ctx context.Context, orderID int64) ([]core_domain.OrderStatusHistory, error) {
			if orderID != 100 {
				t.Fatalf("orderID = %d, want 100", orderID)
			}
			return []core_domain.OrderStatusHistory{{
				ID:        1,
				OrderID:   orderID,
				OldStatus: &oldStatus,
				NewStatus: core_domain.OrderStatusWaitingPickup,
				ChangedBy: 5,
				Comment:   ptr("Ожидает забора"),
				ChangedAt: now,
			}}, nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	history, err := service.GetOrderHistory(context.Background(), 100, 42, core_domain.RoleClient.String())
	if err != nil {
		t.Fatalf("GetOrderHistory returned error: %v", err)
	}
	if len(history) != 1 || history[0].OldStatus == nil || *history[0].OldStatus != core_domain.OrderStatusCreated.String() || history[0].NewStatus != core_domain.OrderStatusWaitingPickup.String() {
		t.Fatalf("unexpected history: %+v", history)
	}
}

func TestUpdateOrderStatusSuccessCreatedToWaitingPickup(t *testing.T) {
	comment := "Назначен забор"
	getCalls := 0

	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			getCalls++
			if getCalls == 1 {
				return makeOrderDetails(orderID, 42, core_domain.HandoverTypePickup, core_domain.OrderStatusCreated), nil
			}
			return makeOrderDetails(orderID, 42, core_domain.HandoverTypePickup, core_domain.OrderStatusWaitingPickup), nil
		},
		updateOrderStatusFn: func(ctx context.Context, orderID int64, status string, changedBy int64, commentArg *string) error {
			if orderID != 100 || status != core_domain.OrderStatusWaitingPickup.String() || changedBy != 5 {
				t.Fatalf("unexpected update args: id=%d status=%q changedBy=%d", orderID, status, changedBy)
			}
			if commentArg == nil || *commentArg != comment {
				t.Fatalf("comment = %v, want %s", commentArg, comment)
			}
			return nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	dto, err := service.UpdateOrderStatus(context.Background(), 100, 5, orders_service.UpdateOrderStatusInput{
		Status:  core_domain.OrderStatusWaitingPickup.String(),
		Comment: &comment,
	})
	if err != nil {
		t.Fatalf("UpdateOrderStatus returned error: %v", err)
	}
	if dto.Status != core_domain.OrderStatusWaitingPickup.String() || getCalls != 2 {
		t.Fatalf("unexpected dto/getCalls: %+v/%d", dto, getCalls)
	}
}

func TestUpdateOrderStatusRejectsInvalidTransition(t *testing.T) {
	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			return makeOrderDetails(orderID, 42, core_domain.HandoverTypePickup, core_domain.OrderStatusCreated), nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	_, err := service.UpdateOrderStatus(context.Background(), 100, 5, orders_service.UpdateOrderStatusInput{Status: core_domain.OrderStatusReceived.String()})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestUpdateOrderStatusRejectsWrongHandoverStatus(t *testing.T) {
	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			return makeOrderDetails(orderID, 42, core_domain.HandoverTypeSelfDelivery, core_domain.OrderStatusCreated), nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	_, err := service.UpdateOrderStatus(context.Background(), 100, 5, orders_service.UpdateOrderStatusInput{Status: core_domain.OrderStatusWaitingPickup.String()})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestCancelOrderSuccess(t *testing.T) {
	comment := "Клиент отменил заявку"
	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			return makeOrderDetails(orderID, 42, core_domain.HandoverTypeSelfDelivery, core_domain.OrderStatusWaitingDelivery), nil
		},
		cancelOrderFn: func(ctx context.Context, orderID int64, changedBy int64, commentArg *string) error {
			if orderID != 100 || changedBy != 42 {
				t.Fatalf("unexpected cancel args: id=%d changedBy=%d", orderID, changedBy)
			}
			if commentArg == nil || *commentArg != comment {
				t.Fatalf("comment = %v, want %s", commentArg, comment)
			}
			return nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	err := service.CancelOrder(context.Background(), 100, 42, core_domain.RoleClient.String(), orders_service.CancelOrderInput{Comment: &comment})
	if err != nil {
		t.Fatalf("CancelOrder returned error: %v", err)
	}
}

func TestCancelOrderRejectsForeignClient(t *testing.T) {
	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			return makeOrderDetails(orderID, 99, core_domain.HandoverTypePickup, core_domain.OrderStatusWaitingPickup), nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	err := service.CancelOrder(context.Background(), 100, 42, core_domain.RoleClient.String(), orders_service.CancelOrderInput{})
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
}

func TestCancelOrderRejectsTerminalOrder(t *testing.T) {
	repo := &fakeOrdersRepository{
		getOrderFn: func(ctx context.Context, orderID int64) (core_domain.OrderDetails, error) {
			return makeOrderDetails(orderID, 42, core_domain.HandoverTypePickup, core_domain.OrderStatusDelivered), nil
		},
	}

	service := orders_service.NewOrdersService(repo)
	err := service.CancelOrder(context.Background(), 100, 42, core_domain.RoleClient.String(), orders_service.CancelOrderInput{})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func validSelfDeliveryInput() orders_service.CreateOrderInput {
	date := "2026-05-20"
	return orders_service.CreateOrderInput{
		ReceivingWarehouseID:   10,
		DestinationWarehouseID: 20,
		ProductTypeID:          30,
		HandoverType:           core_domain.HandoverTypeSelfDelivery.String(),
		SelfDeliveryDate:       &date,
		CargoPlaces: []orders_service.CargoPlaceInput{
			{CargoPlaceTypeID: 7, Quantity: 1},
		},
	}
}
