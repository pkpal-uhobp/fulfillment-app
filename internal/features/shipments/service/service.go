package shipments_service

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type ShipmentsRepository interface {
	CreateShipment(ctx context.Context, shipment core_domain.Shipment) (core_domain.Shipment, error)
	ListShipments(ctx context.Context, filter core_domain.ShipmentFilter) ([]core_domain.Shipment, error)
	GetShipment(ctx context.Context, shipmentID int64) (core_domain.ShipmentDetails, error)
	AddShipmentItem(ctx context.Context, shipmentID int64, cargoItemID int64, changedBy int64, comment *string) (core_domain.ShipmentDetails, error)
	RemoveShipmentItem(ctx context.Context, shipmentID int64, cargoItemID int64) error
	UpdateShipmentStatus(ctx context.Context, shipmentID int64, status string, changedBy int64, comment *string) (core_domain.ShipmentDetails, error)
}

type ShipmentsService struct {
	repo ShipmentsRepository
}

func NewShipmentsService(repo ShipmentsRepository) *ShipmentsService {
	return &ShipmentsService{repo: repo}
}
