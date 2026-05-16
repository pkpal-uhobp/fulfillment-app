package shipments_service

import "context"

func (s *ShipmentsService) AddShipmentItem(
	ctx context.Context,
	shipmentID int64,
	actorID int64,
	actorRole string,
	input AddShipmentItemInput,
) (ShipmentDTO, error) {
	if err := requireLogistOrAdmin(actorRole); err != nil {
		return ShipmentDTO{}, err
	}
	if err := validateShipmentID(shipmentID, "shipment id"); err != nil {
		return ShipmentDTO{}, err
	}
	if err := validateShipmentID(input.CargoItemID, "cargo item id"); err != nil {
		return ShipmentDTO{}, err
	}
	details, err := s.repo.AddShipmentItem(ctx, shipmentID, input.CargoItemID, actorID, input.Comment)
	if err != nil {
		return ShipmentDTO{}, err
	}
	return mapShipmentDetailsToDTO(details), nil
}
