package cargoitems_service

import core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"

func mapCargoItemToDTO(item core_domain.CargoItem) CargoItemDTO {
	return CargoItemDTO{
		ID:                item.ID,
		OrderID:           item.OrderID,
		OrderCargoPlaceID: item.OrderCargoPlaceID,
		CargoPlaceTypeID:  item.CargoPlaceTypeID,
		QRCode:            item.QRCode,
		Status:            item.Status.String(),
		StorageZoneID:     item.StorageZoneID,
		GateID:            item.GateID,
		ReceivedBy:        item.ReceivedBy,
		ShippedBy:         item.ShippedBy,
		ReceivedAt:        item.ReceivedAt,
		ShippedAt:         item.ShippedAt,
		Comment:           item.Comment,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
	}
}

func mapCargoItemsToDTO(items []core_domain.CargoItem) []CargoItemDTO {
	result := make([]CargoItemDTO, 0, len(items))
	for _, item := range items {
		result = append(result, mapCargoItemToDTO(item))
	}
	return result
}

func mapCargoStatusHistoryListToDTO(history []core_domain.CargoStatusHistory) []CargoStatusHistoryDTO {
	result := make([]CargoStatusHistoryDTO, 0, len(history))
	for _, item := range history {
		var oldStatus *string
		if item.OldStatus != nil {
			value := item.OldStatus.String()
			oldStatus = &value
		}

		result = append(result, CargoStatusHistoryDTO{
			ID:          item.ID,
			CargoItemID: item.CargoItemID,
			OldStatus:   oldStatus,
			NewStatus:   item.NewStatus.String(),
			ChangedBy:   item.ChangedBy,
			Comment:     item.Comment,
			ChangedAt:   item.ChangedAt,
		})
	}
	return result
}
