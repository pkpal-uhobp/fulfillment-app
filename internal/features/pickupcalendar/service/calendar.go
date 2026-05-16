package pickupcalendar_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (s *PickupCalendarService) GetCalendar(
	ctx context.Context,
	actorID int64,
	actorRole string,
	filter CalendarFilter,
) ([]PickupCalendarDayDTO, error) {
	if err := validateActorID(actorID); err != nil {
		return nil, err
	}
	if err := requireCalendarReadRole(actorRole); err != nil {
		return nil, err
	}
	if err := validateWarehouseID(filter.WarehouseID); err != nil {
		return nil, err
	}

	dateFrom, dateTo, err := validateDateRange(filter.DateFrom, filter.DateTo)
	if err != nil {
		return nil, err
	}

	days, err := s.repo.ListCalendar(ctx, core_domain.PickupCalendarFilter{
		WarehouseID: filter.WarehouseID,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	})
	if err != nil {
		return nil, fmt.Errorf("list pickup calendar: %w", err)
	}
	return mapCalendarDaysToDTO(days), nil
}
