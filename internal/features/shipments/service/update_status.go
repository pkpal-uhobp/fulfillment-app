package shipments_service

import "context"

func (s *ShipmentsService) UpdateShipmentStatus(
	ctx context.Context,
	shipmentID int64,
	actorID int64,
	actorRole string,
	input UpdateShipmentStatusInput,
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
	nextStatus, err := validateShipmentStatus(input.Status)
	if err != nil {
		return ShipmentDTO{}, err
	}
	if err := validateShipmentStatusTransition(details.Shipment.Status, nextStatus); err != nil {
		return ShipmentDTO{}, err
	}
	updated, err := s.repo.UpdateShipmentStatus(ctx, shipmentID, nextStatus.String(), actorID, input.Comment)
	if err != nil {
		return ShipmentDTO{}, err
	}
	return mapShipmentDetailsToDTO(updated), nil
}
