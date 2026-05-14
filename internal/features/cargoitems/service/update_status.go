package cargoitems_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *CargoItemsService) UpdateCargoItemStatus(
	ctx context.Context,
	cargoItemID int64,
	actorID int64,
	actorRole string,
	input UpdateCargoItemStatusInput,
) (CargoItemDTO, error) {
	if actorRole != core_domain.RoleWorker.String() && actorRole != core_domain.RoleLogist.String() && actorRole != core_domain.RoleAdmin.String() {
		return CargoItemDTO{}, fmt.Errorf("%w: only worker, logist or admin can update cargo item status", core_errors.ErrForbidden)
	}

	current, err := s.getCargoItemWithAccess(ctx, cargoItemID, actorID, actorRole)
	if err != nil {
		return CargoItemDTO{}, err
	}

	nextStatus, err := validateCargoItemStatus(input.Status)
	if err != nil {
		return CargoItemDTO{}, err
	}

	if err := validateCargoStatusTransition(current, nextStatus); err != nil {
		return CargoItemDTO{}, err
	}

	updated, err := s.repo.UpdateCargoItemStatus(
		ctx,
		cargoItemID,
		nextStatus.String(),
		actorID,
		normalizeOptionalString(input.Comment),
	)
	if err != nil {
		return CargoItemDTO{}, err
	}

	return mapCargoItemToDTO(updated), nil
}
