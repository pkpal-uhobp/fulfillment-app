package warehouses_service

import (
	"context"
	"fmt"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *WarehousesService) ActivateWarehouse(ctx context.Context, warehouseID int64) error {
	if warehouseID <= 0 {
		return fmt.Errorf("%w: invalid warehouse id", core_errors.ErrInvalidArgument)
	}

	return s.repo.ActivateWarehouse(ctx, warehouseID)
}
