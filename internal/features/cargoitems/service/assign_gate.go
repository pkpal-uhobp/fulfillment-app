package cargoitems_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *CargoItemsService) AssignGate(
	ctx context.Context,
	cargoItemID int64,
	actorID int64,
	actorRole string,
	input AssignGateInput,
) (CargoItemDTO, error) {
	if actorRole != core_domain.RoleLogist.String() && actorRole != core_domain.RoleAdmin.String() {
		return CargoItemDTO{}, fmt.Errorf("%w: only logist or admin can assign gate", core_errors.ErrForbidden)
	}
	if input.GateID <= 0 {
		return CargoItemDTO{}, fmt.Errorf("%w: gate id is required", core_errors.ErrInvalidArgument)
	}

	current, err := s.getCargoItemWithAccess(ctx, cargoItemID, actorID, actorRole)
	if err != nil {
		return CargoItemDTO{}, err
	}
	if err := ensureCargoItemEditable(current); err != nil {
		return CargoItemDTO{}, err
	}
	if current.Status != core_domain.CargoItemStatusStored && current.Status != core_domain.CargoItemStatusReadyToShip {
		return CargoItemDTO{}, fmt.Errorf("%w: gate can be assigned only to stored or ready_to_ship cargo item", core_errors.ErrConflict)
	}

	ok, err := s.repo.GateBelongsToCargoOrder(ctx, cargoItemID, input.GateID)
	if err != nil {
		return CargoItemDTO{}, err
	}
	if !ok {
		return CargoItemDTO{}, fmt.Errorf("%w: gate does not belong to cargo receiving warehouse", core_errors.ErrInvalidArgument)
	}

	updated, err := s.repo.AssignGate(
		ctx,
		cargoItemID,
		input.GateID,
		actorID,
		normalizeOptionalString(input.Comment),
	)
	if err != nil {
		return CargoItemDTO{}, err
	}

	return mapCargoItemToDTO(updated), nil
}
