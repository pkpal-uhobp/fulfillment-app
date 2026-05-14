package cargoitems_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
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

	filter := core_domain.CargoItemFilter{
		QRCode: qrCode,
		Page:   1,
		Limit:  1,
	}
	if actorRole == core_domain.RoleClient.String() {
		filter.ClientID = &actorID
	}

	items, err := s.repo.ListCargoItems(ctx, filter)
	if err != nil {
		return CargoItemDTO{}, err
	}
	if len(items) == 0 {
		return CargoItemDTO{}, fmt.Errorf("%w: cargo item with qr_code not found", core_errors.ErrNotFound)
	}

	return mapCargoItemToDTO(items[0]), nil
}
