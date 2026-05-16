package pickupcalendar_service_tests

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type fakePickupCalendarRepository struct {
	listCalendarFn func(context.Context, core_domain.PickupCalendarFilter) ([]core_domain.PickupCalendarDay, error)
	createBlockFn  func(context.Context, core_domain.PickupCalendarBlock) (core_domain.PickupCalendarBlock, error)
	deleteBlockFn  func(context.Context, int64) error
	setCapacityFn  func(context.Context, core_domain.PickupCapacity) (core_domain.PickupCapacity, error)
}

func (f *fakePickupCalendarRepository) ListCalendar(ctx context.Context, filter core_domain.PickupCalendarFilter) ([]core_domain.PickupCalendarDay, error) {
	if f.listCalendarFn == nil {
		return nil, fmt.Errorf("unexpected ListCalendar call")
	}
	return f.listCalendarFn(ctx, filter)
}

func (f *fakePickupCalendarRepository) CreateBlock(ctx context.Context, block core_domain.PickupCalendarBlock) (core_domain.PickupCalendarBlock, error) {
	if f.createBlockFn == nil {
		return core_domain.PickupCalendarBlock{}, fmt.Errorf("unexpected CreateBlock call")
	}
	return f.createBlockFn(ctx, block)
}

func (f *fakePickupCalendarRepository) DeleteBlock(ctx context.Context, blockID int64) error {
	if f.deleteBlockFn == nil {
		return fmt.Errorf("unexpected DeleteBlock call")
	}
	return f.deleteBlockFn(ctx, blockID)
}

func (f *fakePickupCalendarRepository) SetCapacity(ctx context.Context, capacity core_domain.PickupCapacity) (core_domain.PickupCapacity, error) {
	if f.setCapacityFn == nil {
		return core_domain.PickupCapacity{}, fmt.Errorf("unexpected SetCapacity call")
	}
	return f.setCapacityFn(ctx, capacity)
}

func ptr[T any](value T) *T {
	return &value
}
