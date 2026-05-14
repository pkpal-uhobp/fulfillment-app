package warehouses_service

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *WarehousesService) ActivateGate(ctx context.Context, gateID int64) error {
	if gateID <= 0 {
		return fmt.Errorf("%w: invalid gate id", core_errors.ErrInvalidArgument)
	}

	return s.repo.ActivateGate(ctx, gateID)
}

func (s *WarehousesService) DeactivateGate(ctx context.Context, gateID int64) error {
	if gateID <= 0 {
		return fmt.Errorf("%w: invalid gate id", core_errors.ErrInvalidArgument)
	}

	return s.repo.DeactivateGate(ctx, gateID)
}
