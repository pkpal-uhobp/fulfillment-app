package shipments_service

import "context"

func (s *ShipmentsService) RemoveShipmentItem(
	ctx context.Context,
	shipmentID int64,
	cargoItemID int64,
	actorID int64,
	actorRole string,
) error {
	if err := requireLogistOrAdmin(actorRole); err != nil {
		return err
	}
	if err := validateShipmentID(shipmentID, "shipment id"); err != nil {
		return err
	}
	if err := validateShipmentID(cargoItemID, "cargo item id"); err != nil {
		return err
	}
	return s.repo.RemoveShipmentItem(ctx, shipmentID, cargoItemID)
}
