package pickupcalendar_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *PickupCalendarService) SetCapacity(
	ctx context.Context,
	actorID int64,
	actorRole string,
	input SetCapacityInput,
) (PickupCapacityDTO, error) {
	if err := validateActorID(actorID); err != nil {
		return PickupCapacityDTO{}, err
	}
	if err := requireCalendarManageRole(actorRole); err != nil {
		return PickupCapacityDTO{}, err
	}
	if err := validateWarehouseID(input.WarehouseID); err != nil {
		return PickupCapacityDTO{}, err
	}
	pickupDate := strings.TrimSpace(input.PickupDate)
	if err := validateDate(pickupDate, "pickup_date"); err != nil {
		return PickupCapacityDTO{}, err
	}
	if input.MaxOrders < 0 {
		return PickupCapacityDTO{}, fmt.Errorf("%w: max_orders cannot be negative", core_errors.ErrInvalidArgument)
	}
	if input.CurrentOrders < 0 {
		return PickupCapacityDTO{}, fmt.Errorf("%w: current_orders cannot be negative", core_errors.ErrInvalidArgument)
	}
	if input.CurrentOrders > input.MaxOrders {
		return PickupCapacityDTO{}, fmt.Errorf("%w: current_orders cannot be greater than max_orders", core_errors.ErrInvalidArgument)
	}

	capacity, err := s.repo.SetCapacity(ctx, core_domain.PickupCapacity{
		WarehouseID:   input.WarehouseID,
		PickupDate:    pickupDate,
		MaxOrders:     input.MaxOrders,
		CurrentOrders: input.CurrentOrders,
		IsClosed:      input.IsClosed,
	})
	if err != nil {
		return PickupCapacityDTO{}, fmt.Errorf("set pickup capacity: %w", err)
	}
	return mapCapacityToDTO(capacity), nil
}
