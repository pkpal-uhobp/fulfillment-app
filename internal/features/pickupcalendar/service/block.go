package pickupcalendar_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (s *PickupCalendarService) BlockDate(
	ctx context.Context,
	actorID int64,
	actorRole string,
	input BlockDateInput,
) (PickupCalendarBlockDTO, error) {
	if err := validateActorID(actorID); err != nil {
		return PickupCalendarBlockDTO{}, err
	}
	if err := requireCalendarManageRole(actorRole); err != nil {
		return PickupCalendarBlockDTO{}, err
	}
	if err := validateWarehouseID(input.WarehouseID); err != nil {
		return PickupCalendarBlockDTO{}, err
	}
	blockedDate := strings.TrimSpace(input.BlockedDate)
	if err := validateDate(blockedDate, "blocked_date"); err != nil {
		return PickupCalendarBlockDTO{}, err
	}

	created, err := s.repo.CreateBlock(ctx, core_domain.PickupCalendarBlock{
		WarehouseID: input.WarehouseID,
		BlockedDate: blockedDate,
		Reason:      normalizeOptionalString(input.Reason),
		CreatedBy:   actorID,
	})
	if err != nil {
		return PickupCalendarBlockDTO{}, fmt.Errorf("create pickup calendar block: %w", err)
	}
	return mapBlockToDTO(created), nil
}

func (s *PickupCalendarService) UnblockDate(
	ctx context.Context,
	actorID int64,
	actorRole string,
	blockID int64,
) error {
	if err := validateActorID(actorID); err != nil {
		return err
	}
	if err := requireCalendarManageRole(actorRole); err != nil {
		return err
	}
	if err := validateBlockID(blockID); err != nil {
		return err
	}
	if err := s.repo.DeleteBlock(ctx, blockID); err != nil {
		return fmt.Errorf("delete pickup calendar block: %w", err)
	}
	return nil
}
