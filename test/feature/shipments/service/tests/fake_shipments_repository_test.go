package shipments_service_tests

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type fakeShipmentsRepository struct {
	createShipmentFn       func(context.Context, core_domain.Shipment) (core_domain.Shipment, error)
	listShipmentsFn        func(context.Context, core_domain.ShipmentFilter) ([]core_domain.Shipment, error)
	getShipmentFn          func(context.Context, int64) (core_domain.ShipmentDetails, error)
	addShipmentItemFn      func(context.Context, int64, int64, int64, *string) (core_domain.ShipmentDetails, error)
	removeShipmentItemFn   func(context.Context, int64, int64) error
	updateShipmentStatusFn func(context.Context, int64, string, int64, *string) (core_domain.ShipmentDetails, error)
}

func (f *fakeShipmentsRepository) CreateShipment(ctx context.Context, shipment core_domain.Shipment) (core_domain.Shipment, error) {
	if f.createShipmentFn == nil {
		return core_domain.Shipment{}, fmt.Errorf("unexpected CreateShipment call")
	}
	return f.createShipmentFn(ctx, shipment)
}

func (f *fakeShipmentsRepository) ListShipments(ctx context.Context, filter core_domain.ShipmentFilter) ([]core_domain.Shipment, error) {
	if f.listShipmentsFn == nil {
		return nil, fmt.Errorf("unexpected ListShipments call")
	}
	return f.listShipmentsFn(ctx, filter)
}

func (f *fakeShipmentsRepository) GetShipment(ctx context.Context, shipmentID int64) (core_domain.ShipmentDetails, error) {
	if f.getShipmentFn == nil {
		return core_domain.ShipmentDetails{}, fmt.Errorf("unexpected GetShipment call")
	}
	return f.getShipmentFn(ctx, shipmentID)
}

func (f *fakeShipmentsRepository) AddShipmentItem(ctx context.Context, shipmentID int64, cargoItemID int64, changedBy int64, comment *string) (core_domain.ShipmentDetails, error) {
	if f.addShipmentItemFn == nil {
		return core_domain.ShipmentDetails{}, fmt.Errorf("unexpected AddShipmentItem call")
	}
	return f.addShipmentItemFn(ctx, shipmentID, cargoItemID, changedBy, comment)
}

func (f *fakeShipmentsRepository) RemoveShipmentItem(ctx context.Context, shipmentID int64, cargoItemID int64) error {
	if f.removeShipmentItemFn == nil {
		return fmt.Errorf("unexpected RemoveShipmentItem call")
	}
	return f.removeShipmentItemFn(ctx, shipmentID, cargoItemID)
}

func (f *fakeShipmentsRepository) UpdateShipmentStatus(ctx context.Context, shipmentID int64, status string, changedBy int64, comment *string) (core_domain.ShipmentDetails, error) {
	if f.updateShipmentStatusFn == nil {
		return core_domain.ShipmentDetails{}, fmt.Errorf("unexpected UpdateShipmentStatus call")
	}
	return f.updateShipmentStatusFn(ctx, shipmentID, status, changedBy, comment)
}

func ptr[T any](value T) *T { return &value }
