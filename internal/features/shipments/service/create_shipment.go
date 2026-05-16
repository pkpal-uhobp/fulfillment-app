package shipments_service

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (s *ShipmentsService) CreateShipment(
	ctx context.Context,
	actorID int64,
	actorRole string,
	input CreateShipmentInput,
) (ShipmentDTO, error) {
	if err := requireLogistOrAdmin(actorRole); err != nil {
		return ShipmentDTO{}, err
	}
	if err := validateShipmentID(input.DestinationWarehouseID, "destination warehouse id"); err != nil {
		return ShipmentDTO{}, err
	}
	if err := validateShipmentID(input.GateID, "gate id"); err != nil {
		return ShipmentDTO{}, err
	}
	plannedAt, err := parsePlannedDeparture(input.PlannedDepartureAt)
	if err != nil {
		return ShipmentDTO{}, err
	}

	shipment := core_domain.Shipment{
		DestinationWarehouseID: input.DestinationWarehouseID,
		GateID:                 input.GateID,
		PlannedDepartureAt:     plannedAt,
		Status:                 core_domain.ShipmentStatusPlanned,
		CreatedBy:              actorID,
	}
	created, err := s.repo.CreateShipment(ctx, shipment)
	if err != nil {
		return ShipmentDTO{}, err
	}
	return mapShipmentToDTO(created), nil
}
