package shipments_service_tests

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	shipments_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/service"
)

func TestCreateShipmentSuccess(t *testing.T) {
	planned := "2026-06-20T12:00:00Z"
	repo := &fakeShipmentsRepository{
		createShipmentFn: func(_ context.Context, shipment core_domain.Shipment) (core_domain.Shipment, error) {
			if shipment.DestinationWarehouseID != 10 || shipment.GateID != 20 || shipment.CreatedBy != 7 {
				t.Fatalf("unexpected shipment: %+v", shipment)
			}
			if shipment.Status != core_domain.ShipmentStatusPlanned {
				t.Fatalf("status = %s, want planned", shipment.Status)
			}
			shipment.ID = 100
			shipment.CreatedAt = time.Now()
			return shipment, nil
		},
	}
	service := shipments_service.NewShipmentsService(repo)

	created, err := service.CreateShipment(context.Background(), 7, core_domain.RoleLogist.String(), shipments_service.CreateShipmentInput{
		DestinationWarehouseID: 10,
		GateID:                 20,
		PlannedDepartureAt:     planned,
	})
	if err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	if created.ID != 100 || created.Status != core_domain.ShipmentStatusPlanned.String() {
		t.Fatalf("unexpected dto: %+v", created)
	}
}

func TestCreateShipmentRejectsClient(t *testing.T) {
	service := shipments_service.NewShipmentsService(&fakeShipmentsRepository{})
	_, err := service.CreateShipment(context.Background(), 1, core_domain.RoleClient.String(), shipments_service.CreateShipmentInput{
		DestinationWarehouseID: 10,
		GateID:                 20,
		PlannedDepartureAt:     "2026-06-20T12:00:00Z",
	})
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("error = %v, want forbidden", err)
	}
}

func TestCreateShipmentRejectsInvalidDate(t *testing.T) {
	service := shipments_service.NewShipmentsService(&fakeShipmentsRepository{})
	_, err := service.CreateShipment(context.Background(), 1, core_domain.RoleAdmin.String(), shipments_service.CreateShipmentInput{
		DestinationWarehouseID: 10,
		GateID:                 20,
		PlannedDepartureAt:     "20.06.2026",
	})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("error = %v, want invalid argument", err)
	}
}

func TestAddShipmentItemSuccess(t *testing.T) {
	comment := "add to shipment"
	repo := &fakeShipmentsRepository{
		addShipmentItemFn: func(_ context.Context, shipmentID int64, cargoItemID int64, changedBy int64, gotComment *string) (core_domain.ShipmentDetails, error) {
			if shipmentID != 100 || cargoItemID != 55 || changedBy != 7 {
				t.Fatalf("unexpected args: shipment=%d cargo=%d changedBy=%d", shipmentID, cargoItemID, changedBy)
			}
			if gotComment == nil || *gotComment != comment {
				t.Fatalf("comment = %v, want %q", gotComment, comment)
			}
			return core_domain.ShipmentDetails{
				Shipment: core_domain.Shipment{ID: shipmentID, Status: core_domain.ShipmentStatusPlanned},
				Items: []core_domain.ShipmentItemDetails{{
					Item:      core_domain.ShipmentItem{ID: 1, ShipmentID: shipmentID, CargoItemID: cargoItemID},
					CargoItem: core_domain.CargoItem{ID: cargoItemID, OrderID: 9, QRCode: "QR-55", Status: core_domain.CargoItemStatusReadyToShip},
				}},
			}, nil
		},
	}
	service := shipments_service.NewShipmentsService(repo)

	shipment, err := service.AddShipmentItem(context.Background(), 100, 7, core_domain.RoleLogist.String(), shipments_service.AddShipmentItemInput{
		CargoItemID: 55,
		Comment:     &comment,
	})
	if err != nil {
		t.Fatalf("AddShipmentItem() error = %v", err)
	}
	if len(shipment.Items) != 1 || shipment.Items[0].CargoItemID != 55 {
		t.Fatalf("unexpected shipment items: %+v", shipment.Items)
	}
}

func TestUpdateShipmentStatusSuccessLoadingToShipped(t *testing.T) {
	called := false
	repo := &fakeShipmentsRepository{
		getShipmentFn: func(_ context.Context, shipmentID int64) (core_domain.ShipmentDetails, error) {
			if shipmentID != 100 {
				t.Fatalf("shipmentID = %d, want 100", shipmentID)
			}
			return core_domain.ShipmentDetails{Shipment: core_domain.Shipment{ID: 100, Status: core_domain.ShipmentStatusLoading}}, nil
		},
		updateShipmentStatusFn: func(_ context.Context, shipmentID int64, status string, changedBy int64, comment *string) (core_domain.ShipmentDetails, error) {
			if shipmentID != 100 || status != core_domain.ShipmentStatusShipped.String() || changedBy != 7 {
				t.Fatalf("unexpected args: shipment=%d status=%s changedBy=%d", shipmentID, status, changedBy)
			}
			called = true
			return core_domain.ShipmentDetails{Shipment: core_domain.Shipment{ID: 100, Status: core_domain.ShipmentStatusShipped}}, nil
		},
	}
	service := shipments_service.NewShipmentsService(repo)

	updated, err := service.UpdateShipmentStatus(context.Background(), 100, 7, core_domain.RoleAdmin.String(), shipments_service.UpdateShipmentStatusInput{
		Status: core_domain.ShipmentStatusShipped.String(),
	})
	if err != nil {
		t.Fatalf("UpdateShipmentStatus() error = %v", err)
	}
	if !called {
		t.Fatal("UpdateShipmentStatus repository was not called")
	}
	if updated.Status != core_domain.ShipmentStatusShipped.String() {
		t.Fatalf("status = %s, want shipped", updated.Status)
	}
}

func TestUpdateShipmentStatusRejectsInvalidTransition(t *testing.T) {
	repo := &fakeShipmentsRepository{
		getShipmentFn: func(_ context.Context, shipmentID int64) (core_domain.ShipmentDetails, error) {
			return core_domain.ShipmentDetails{Shipment: core_domain.Shipment{ID: shipmentID, Status: core_domain.ShipmentStatusPlanned}}, nil
		},
	}
	service := shipments_service.NewShipmentsService(repo)

	_, err := service.UpdateShipmentStatus(context.Background(), 100, 7, core_domain.RoleLogist.String(), shipments_service.UpdateShipmentStatusInput{
		Status: core_domain.ShipmentStatusCompleted.String(),
	})
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("error = %v, want conflict", err)
	}
}

func TestRemoveShipmentItemSuccess(t *testing.T) {
	called := false
	repo := &fakeShipmentsRepository{
		removeShipmentItemFn: func(_ context.Context, shipmentID int64, cargoItemID int64) error {
			if shipmentID != 100 || cargoItemID != 55 {
				t.Fatalf("unexpected ids: shipment=%d cargo=%d", shipmentID, cargoItemID)
			}
			called = true
			return nil
		},
	}
	service := shipments_service.NewShipmentsService(repo)
	if err := service.RemoveShipmentItem(context.Background(), 100, 55, 7, core_domain.RoleLogist.String()); err != nil {
		t.Fatalf("RemoveShipmentItem() error = %v", err)
	}
	if !called {
		t.Fatal("RemoveShipmentItem repository was not called")
	}
}
