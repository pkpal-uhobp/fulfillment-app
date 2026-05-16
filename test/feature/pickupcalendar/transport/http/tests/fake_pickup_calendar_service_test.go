package pickupcalendar_transport_http_tests

import (
	"context"
	"fmt"
	"testing"

	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	pickupcalendar_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/service"
	pickupcalendar_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/transport/http"
)

type fakePickupCalendarService struct {
	getCalendarFn func(context.Context, int64, string, pickupcalendar_service.CalendarFilter) ([]pickupcalendar_service.PickupCalendarDayDTO, error)
	blockDateFn   func(context.Context, int64, string, pickupcalendar_service.BlockDateInput) (pickupcalendar_service.PickupCalendarBlockDTO, error)
	unblockDateFn func(context.Context, int64, string, int64) error
	setCapacityFn func(context.Context, int64, string, pickupcalendar_service.SetCapacityInput) (pickupcalendar_service.PickupCapacityDTO, error)
}

func (f *fakePickupCalendarService) GetCalendar(ctx context.Context, actorID int64, actorRole string, filter pickupcalendar_service.CalendarFilter) ([]pickupcalendar_service.PickupCalendarDayDTO, error) {
	if f.getCalendarFn == nil {
		return nil, fmt.Errorf("unexpected GetCalendar call")
	}
	return f.getCalendarFn(ctx, actorID, actorRole, filter)
}

func (f *fakePickupCalendarService) BlockDate(ctx context.Context, actorID int64, actorRole string, input pickupcalendar_service.BlockDateInput) (pickupcalendar_service.PickupCalendarBlockDTO, error) {
	if f.blockDateFn == nil {
		return pickupcalendar_service.PickupCalendarBlockDTO{}, fmt.Errorf("unexpected BlockDate call")
	}
	return f.blockDateFn(ctx, actorID, actorRole, input)
}

func (f *fakePickupCalendarService) UnblockDate(ctx context.Context, actorID int64, actorRole string, blockID int64) error {
	if f.unblockDateFn == nil {
		return fmt.Errorf("unexpected UnblockDate call")
	}
	return f.unblockDateFn(ctx, actorID, actorRole, blockID)
}

func (f *fakePickupCalendarService) SetCapacity(ctx context.Context, actorID int64, actorRole string, input pickupcalendar_service.SetCapacityInput) (pickupcalendar_service.PickupCapacityDTO, error) {
	if f.setCapacityFn == nil {
		return pickupcalendar_service.PickupCapacityDTO{}, fmt.Errorf("unexpected SetCapacity call")
	}
	return f.setCapacityFn(ctx, actorID, actorRole, input)
}

func newTestHandler(t *testing.T, service *fakePickupCalendarService) *pickupcalendar_http.PickupCalendarHTTPHandler {
	t.Helper()
	log, err := core_logger.NewLogger(core_logger.LoggerConfig{
		Level:  "debug",
		Folder: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	t.Cleanup(log.Close)
	return pickupcalendar_http.NewPickupCalendarHTTPHandler(log, service)
}

func ptr[T any](value T) *T {
	return &value
}
