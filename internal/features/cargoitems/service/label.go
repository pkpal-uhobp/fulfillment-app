package cargoitems_service

import (
	"context"
	"fmt"
)

func (s *CargoItemsService) GetCargoItemLabel(
	ctx context.Context,
	cargoItemID int64,
	actorID int64,
	actorRole string,
) (CargoItemLabelDTO, error) {
	item, err := s.getCargoItemWithAccess(ctx, cargoItemID, actorID, actorRole)
	if err != nil {
		return CargoItemLabelDTO{}, err
	}

	return CargoItemLabelDTO{
		CargoItemID:       item.ID,
		OrderID:           item.OrderID,
		OrderCargoPlaceID: item.OrderCargoPlaceID,
		CargoPlaceTypeID:  item.CargoPlaceTypeID,
		QRCode:            item.QRCode,
		QRCodeValue:       item.QRCode,
		Status:            item.Status.String(),
		StorageZoneID:     item.StorageZoneID,
		GateID:            item.GateID,
		ReceivedAt:        item.ReceivedAt,
		ShippedAt:         item.ShippedAt,
		LabelText:         fmt.Sprintf("ORDER-%d | CARGO-%d | %s", item.OrderID, item.ID, item.QRCode),
	}, nil
}
