package shipments_service

import core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"

func mapShipmentToDTO(shipment core_domain.Shipment) ShipmentDTO {
	return ShipmentDTO{
		ID:                     shipment.ID,
		DestinationWarehouseID: shipment.DestinationWarehouseID,
		GateID:                 shipment.GateID,
		PlannedDepartureAt:     shipment.PlannedDepartureAt,
		ActualDepartureAt:      shipment.ActualDepartureAt,
		Status:                 shipment.Status.String(),
		CreatedBy:              shipment.CreatedBy,
		CreatedAt:              shipment.CreatedAt,
	}
}

func mapShipmentDetailsToDTO(details core_domain.ShipmentDetails) ShipmentDTO {
	dto := mapShipmentToDTO(details.Shipment)
	dto.Items = make([]ShipmentItemDTO, 0, len(details.Items))
	for _, item := range details.Items {
		dto.Items = append(dto.Items, mapShipmentItemToDTO(item))
	}
	return dto
}

func mapShipmentItemToDTO(item core_domain.ShipmentItemDetails) ShipmentItemDTO {
	return ShipmentItemDTO{
		ID:            item.Item.ID,
		ShipmentID:    item.Item.ShipmentID,
		CargoItemID:   item.Item.CargoItemID,
		OrderID:       item.CargoItem.OrderID,
		QRCode:        item.CargoItem.QRCode,
		Status:        item.CargoItem.Status.String(),
		StorageZoneID: item.CargoItem.StorageZoneID,
		GateID:        item.CargoItem.GateID,
		ReceivedAt:    item.CargoItem.ReceivedAt,
		ShippedAt:     item.CargoItem.ShippedAt,
	}
}

func mapShipmentsToDTO(shipments []core_domain.Shipment) []ShipmentDTO {
	result := make([]ShipmentDTO, 0, len(shipments))
	for _, shipment := range shipments {
		result = append(result, mapShipmentToDTO(shipment))
	}
	return result
}
