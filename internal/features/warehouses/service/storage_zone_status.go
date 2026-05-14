package warehouses_service

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *WarehousesService) ActivateStorageZone(ctx context.Context, zoneID int64) error {
	if zoneID <= 0 {
		return fmt.Errorf("%w: invalid storage zone id", core_errors.ErrInvalidArgument)
	}

	return s.repo.ActivateStorageZone(ctx, zoneID)
}

func (s *WarehousesService) DeactivateStorageZone(ctx context.Context, zoneID int64) error {
	if zoneID <= 0 {
		return fmt.Errorf("%w: invalid storage zone id", core_errors.ErrInvalidArgument)
	}

	return s.repo.DeactivateStorageZone(ctx, zoneID)
}
