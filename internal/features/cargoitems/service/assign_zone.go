package cargoitems_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *CargoItemsService) AssignStorageZone(
	ctx context.Context,
	cargoItemID int64,
	actorID int64,
	actorRole string,
	input AssignStorageZoneInput,
) (CargoItemDTO, error) {
	if actorRole != core_domain.RoleLogist.String() && actorRole != core_domain.RoleAdmin.String() {
		return CargoItemDTO{}, fmt.Errorf("%w: only logist or admin can assign storage zone", core_errors.ErrForbidden)
	}
	if input.StorageZoneID <= 0 {
		return CargoItemDTO{}, fmt.Errorf("%w: storage zone id is required", core_errors.ErrInvalidArgument)
	}

	current, err := s.getCargoItemWithAccess(ctx, cargoItemID, actorID, actorRole)
	if err != nil {
		return CargoItemDTO{}, err
	}
	if err := ensureCargoItemEditable(current); err != nil {
		return CargoItemDTO{}, err
	}
	if current.Status != core_domain.CargoItemStatusAccepted && current.Status != core_domain.CargoItemStatusStored {
		return CargoItemDTO{}, fmt.Errorf("%w: storage zone can be assigned only to accepted or stored cargo item", core_errors.ErrConflict)
	}

	ok, err := s.repo.StorageZoneBelongsToCargoOrder(ctx, cargoItemID, input.StorageZoneID)
	if err != nil {
		return CargoItemDTO{}, err
	}
	if !ok {
		return CargoItemDTO{}, fmt.Errorf("%w: storage zone does not belong to cargo receiving warehouse", core_errors.ErrInvalidArgument)
	}

	updated, err := s.repo.AssignStorageZone(
		ctx,
		cargoItemID,
		input.StorageZoneID,
		actorID,
		normalizeOptionalString(input.Comment),
	)
	if err != nil {
		return CargoItemDTO{}, err
	}

	return mapCargoItemToDTO(updated), nil
}
