package shipments_service

import "context"

func (s *ShipmentsService) GetShipment(
	ctx context.Context,
	shipmentID int64,
	actorID int64,
	actorRole string,
) (ShipmentDTO, error) {
	if err := requireLogistOrAdmin(actorRole); err != nil {
		return ShipmentDTO{}, err
	}
	if err := validateShipmentID(shipmentID, "shipment id"); err != nil {
		return ShipmentDTO{}, err
	}
	details, err := s.repo.GetShipment(ctx, shipmentID)
	if err != nil {
		return ShipmentDTO{}, err
	}
	return mapShipmentDetailsToDTO(details), nil
}
