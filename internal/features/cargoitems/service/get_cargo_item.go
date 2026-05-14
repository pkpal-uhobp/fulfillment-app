package cargoitems_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *CargoItemsService) GetCargoItem(
	ctx context.Context,
	cargoItemID int64,
	actorID int64,
	actorRole string,
) (CargoItemDTO, error) {
	item, err := s.getCargoItemWithAccess(ctx, cargoItemID, actorID, actorRole)
	if err != nil {
		return CargoItemDTO{}, err
	}

	return mapCargoItemToDTO(item), nil
}

func (s *CargoItemsService) getCargoItemWithAccess(
	ctx context.Context,
	cargoItemID int64,
	actorID int64,
	actorRole string,
) (core_domain.CargoItem, error) {
	if cargoItemID <= 0 {
		return core_domain.CargoItem{}, fmt.Errorf("%w: invalid cargo item id", core_errors.ErrInvalidArgument)
	}

	item, err := s.repo.GetCargoItem(ctx, cargoItemID)
	if err != nil {
		return core_domain.CargoItem{}, err
	}

	if actorRole == core_domain.RoleClient.String() {
		owns, err := s.repo.ClientOwnsCargoItem(ctx, cargoItemID, actorID)
		if err != nil {
			return core_domain.CargoItem{}, err
		}
		if !owns {
			return core_domain.CargoItem{}, fmt.Errorf("%w: cargo item does not belong to client", core_errors.ErrForbidden)
		}
	}

	return item, nil
}
