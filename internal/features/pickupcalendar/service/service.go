package pickupcalendar_service

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type PickupCalendarRepository interface {
	ListCalendar(
		ctx context.Context,
		filter core_domain.PickupCalendarFilter,
	) ([]core_domain.PickupCalendarDay, error)
	CreateBlock(
		ctx context.Context,
		block core_domain.PickupCalendarBlock,
	) (core_domain.PickupCalendarBlock, error)
	DeleteBlock(
		ctx context.Context,
		blockID int64,
	) error
	SetCapacity(
		ctx context.Context,
		capacity core_domain.PickupCapacity,
	) (core_domain.PickupCapacity, error)
}

type PickupCalendarService struct {
	repo PickupCalendarRepository
}

func NewPickupCalendarService(repo PickupCalendarRepository) *PickupCalendarService {
	return &PickupCalendarService{repo: repo}
}
