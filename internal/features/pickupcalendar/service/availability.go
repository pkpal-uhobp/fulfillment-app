package pickupcalendar_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *PickupCalendarService) EnsureDateAvailable(
	ctx context.Context,
	warehouseID int64,
	date string,
) error {
	if err := validateWarehouseID(warehouseID); err != nil {
		return err
	}

	pickupDate := strings.TrimSpace(date)
	if err := validateDate(pickupDate, "pickup_date"); err != nil {
		return err
	}

	days, err := s.repo.ListCalendar(ctx, core_domain.PickupCalendarFilter{
		WarehouseID: warehouseID,
		DateFrom:    pickupDate,
		DateTo:      pickupDate,
	})
	if err != nil {
		return fmt.Errorf("list pickup calendar: %w", err)
	}
	if len(days) == 0 {
		return nil
	}
	if days[0].IsClosed {
		return fmt.Errorf("%w: pickup calendar date is closed", core_errors.ErrConflict)
	}
	return nil
}
