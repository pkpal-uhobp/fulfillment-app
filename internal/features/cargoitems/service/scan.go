package cargoitems_service

import (
	"context"
	"fmt"
	"strings"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *CargoItemsService) ScanCargoItem(
	ctx context.Context,
	actorID int64,
	actorRole string,
	qrCode string,
) (CargoItemDTO, error) {
	qrCode = strings.TrimSpace(qrCode)
	if qrCode == "" {
		return CargoItemDTO{}, fmt.Errorf("%w: qr_code is required", core_errors.ErrInvalidArgument)
	}

	// Используем уже существующий ListCargoItems, потому что в нём уже есть
	// ролевая фильтрация: клиент видит только грузовые места своих заявок,
	// а рабочий, логист и администратор могут искать по QR в рамках своих прав.
	items, err := s.ListCargoItems(ctx, actorID, actorRole, CargoItemFilter{
		QRCode: qrCode,
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return CargoItemDTO{}, err
	}

	if len(items) == 0 {
		return CargoItemDTO{}, fmt.Errorf("%w: cargo item with qr_code not found", core_errors.ErrNotFound)
	}

	return items[0], nil
}
