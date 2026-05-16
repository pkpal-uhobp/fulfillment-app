package shipments_transport_http_tests

import (
	"context"
	"fmt"
	"testing"

	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	shipments_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/service"
	shipments_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/transport/http"
)

type fakeShipmentsService struct {
	createShipmentFn       func(context.Context, int64, string, shipments_service.CreateShipmentInput) (shipments_service.ShipmentDTO, error)
	listShipmentsFn        func(context.Context, int64, string, shipments_service.ShipmentFilter) ([]shipments_service.ShipmentDTO, error)
	getShipmentFn          func(context.Context, int64, int64, string) (shipments_service.ShipmentDTO, error)
	addShipmentItemFn      func(context.Context, int64, int64, string, shipments_service.AddShipmentItemInput) (shipments_service.ShipmentDTO, error)
	removeShipmentItemFn   func(context.Context, int64, int64, int64, string) error
	updateShipmentStatusFn func(context.Context, int64, int64, string, shipments_service.UpdateShipmentStatusInput) (shipments_service.ShipmentDTO, error)
}

func (f *fakeShipmentsService) CreateShipment(ctx context.Context, actorID int64, actorRole string, input shipments_service.CreateShipmentInput) (shipments_service.ShipmentDTO, error) {
	if f.createShipmentFn == nil {
		return shipments_service.ShipmentDTO{}, fmt.Errorf("unexpected CreateShipment call")
	}
	return f.createShipmentFn(ctx, actorID, actorRole, input)
}

func (f *fakeShipmentsService) ListShipments(ctx context.Context, actorID int64, actorRole string, filter shipments_service.ShipmentFilter) ([]shipments_service.ShipmentDTO, error) {
	if f.listShipmentsFn == nil {
		return nil, fmt.Errorf("unexpected ListShipments call")
	}
	return f.listShipmentsFn(ctx, actorID, actorRole, filter)
}

func (f *fakeShipmentsService) GetShipment(ctx context.Context, shipmentID int64, actorID int64, actorRole string) (shipments_service.ShipmentDTO, error) {
	if f.getShipmentFn == nil {
		return shipments_service.ShipmentDTO{}, fmt.Errorf("unexpected GetShipment call")
	}
	return f.getShipmentFn(ctx, shipmentID, actorID, actorRole)
}

func (f *fakeShipmentsService) AddShipmentItem(ctx context.Context, shipmentID int64, actorID int64, actorRole string, input shipments_service.AddShipmentItemInput) (shipments_service.ShipmentDTO, error) {
	if f.addShipmentItemFn == nil {
		return shipments_service.ShipmentDTO{}, fmt.Errorf("unexpected AddShipmentItem call")
	}
	return f.addShipmentItemFn(ctx, shipmentID, actorID, actorRole, input)
}

func (f *fakeShipmentsService) RemoveShipmentItem(ctx context.Context, shipmentID int64, cargoItemID int64, actorID int64, actorRole string) error {
	if f.removeShipmentItemFn == nil {
		return fmt.Errorf("unexpected RemoveShipmentItem call")
	}
	return f.removeShipmentItemFn(ctx, shipmentID, cargoItemID, actorID, actorRole)
}

func (f *fakeShipmentsService) UpdateShipmentStatus(ctx context.Context, shipmentID int64, actorID int64, actorRole string, input shipments_service.UpdateShipmentStatusInput) (shipments_service.ShipmentDTO, error) {
	if f.updateShipmentStatusFn == nil {
		return shipments_service.ShipmentDTO{}, fmt.Errorf("unexpected UpdateShipmentStatus call")
	}
	return f.updateShipmentStatusFn(ctx, shipmentID, actorID, actorRole, input)
}

func newTestHandler(t *testing.T, service *fakeShipmentsService) *shipments_http.ShipmentsHTTPHandler {
	t.Helper()
	log, err := core_logger.NewLogger(core_logger.LoggerConfig{
		Level:  "debug",
		Folder: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	t.Cleanup(log.Close)
	return shipments_http.NewShipmentsHTTPHandler(log, service)
}
