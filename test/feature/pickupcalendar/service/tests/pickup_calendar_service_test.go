package pickupcalendar_service_tests

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	pickupcalendar_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/service"
)

func TestGetCalendarSuccess(t *testing.T) {
	repo := &fakePickupCalendarRepository{
		listCalendarFn: func(_ context.Context, filter core_domain.PickupCalendarFilter) ([]core_domain.PickupCalendarDay, error) {
			if filter.WarehouseID != 10 {
				t.Fatalf("WarehouseID = %d, want 10", filter.WarehouseID)
			}
			if filter.DateFrom != "2026-06-01" || filter.DateTo != "2026-06-02" {
				t.Fatalf("date range = %s..%s, want 2026-06-01..2026-06-02", filter.DateFrom, filter.DateTo)
			}
			return []core_domain.PickupCalendarDay{
				{WarehouseID: 10, Date: "2026-06-01", MaxOrders: 5, CurrentOrders: 1, IsClosed: false},
				{WarehouseID: 10, Date: "2026-06-02", MaxOrders: 5, CurrentOrders: 5, IsClosed: true},
			}, nil
		},
	}
	service := pickupcalendar_service.NewPickupCalendarService(repo)

	days, err := service.GetCalendar(context.Background(), 1, core_domain.RoleClient.String(), pickupcalendar_service.CalendarFilter{
		WarehouseID: 10,
		DateFrom:    "2026-06-01",
		DateTo:      "2026-06-02",
	})
	if err != nil {
		t.Fatalf("GetCalendar() error = %v", err)
	}
	if len(days) != 2 {
		t.Fatalf("days len = %d, want 2", len(days))
	}
	if !days[1].IsClosed {
		t.Fatal("second day must be closed")
	}
}

func TestGetCalendarRejectsWorker(t *testing.T) {
	service := pickupcalendar_service.NewPickupCalendarService(&fakePickupCalendarRepository{})
	_, err := service.GetCalendar(context.Background(), 1, core_domain.RoleWorker.String(), pickupcalendar_service.CalendarFilter{WarehouseID: 10})
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("error = %v, want forbidden", err)
	}
}

func TestBlockDateSuccess(t *testing.T) {
	reason := "Перегрузка склада"
	repo := &fakePickupCalendarRepository{
		createBlockFn: func(_ context.Context, block core_domain.PickupCalendarBlock) (core_domain.PickupCalendarBlock, error) {
			if block.WarehouseID != 10 {
				t.Fatalf("WarehouseID = %d, want 10", block.WarehouseID)
			}
			if block.BlockedDate != "2026-06-10" {
				t.Fatalf("BlockedDate = %q, want 2026-06-10", block.BlockedDate)
			}
			if block.CreatedBy != 7 {
				t.Fatalf("CreatedBy = %d, want 7", block.CreatedBy)
			}
			if block.Reason == nil || *block.Reason != reason {
				t.Fatalf("Reason = %v, want %q", block.Reason, reason)
			}
			block.ID = 100
			block.CreatedAt = time.Now()
			return block, nil
		},
	}
	service := pickupcalendar_service.NewPickupCalendarService(repo)

	block, err := service.BlockDate(context.Background(), 7, core_domain.RoleLogist.String(), pickupcalendar_service.BlockDateInput{
		WarehouseID: 10,
		BlockedDate: "2026-06-10",
		Reason:      &reason,
	})
	if err != nil {
		t.Fatalf("BlockDate() error = %v", err)
	}
	if block.ID != 100 {
		t.Fatalf("block.ID = %d, want 100", block.ID)
	}
}

func TestBlockDateRejectsClient(t *testing.T) {
	service := pickupcalendar_service.NewPickupCalendarService(&fakePickupCalendarRepository{})
	_, err := service.BlockDate(context.Background(), 1, core_domain.RoleClient.String(), pickupcalendar_service.BlockDateInput{
		WarehouseID: 10,
		BlockedDate: "2026-06-10",
	})
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("error = %v, want forbidden", err)
	}
}

func TestBlockDateRejectsInvalidDate(t *testing.T) {
	service := pickupcalendar_service.NewPickupCalendarService(&fakePickupCalendarRepository{})
	_, err := service.BlockDate(context.Background(), 1, core_domain.RoleAdmin.String(), pickupcalendar_service.BlockDateInput{
		WarehouseID: 10,
		BlockedDate: "10.06.2026",
	})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("error = %v, want invalid argument", err)
	}
}

func TestUnblockDateSuccess(t *testing.T) {
	called := false
	repo := &fakePickupCalendarRepository{
		deleteBlockFn: func(_ context.Context, blockID int64) error {
			if blockID != 55 {
				t.Fatalf("blockID = %d, want 55", blockID)
			}
			called = true
			return nil
		},
	}
	service := pickupcalendar_service.NewPickupCalendarService(repo)

	if err := service.UnblockDate(context.Background(), 2, core_domain.RoleAdmin.String(), 55); err != nil {
		t.Fatalf("UnblockDate() error = %v", err)
	}
	if !called {
		t.Fatal("DeleteBlock was not called")
	}
}

func TestSetCapacitySuccess(t *testing.T) {
	repo := &fakePickupCalendarRepository{
		setCapacityFn: func(_ context.Context, capacity core_domain.PickupCapacity) (core_domain.PickupCapacity, error) {
			if capacity.WarehouseID != 10 {
				t.Fatalf("WarehouseID = %d, want 10", capacity.WarehouseID)
			}
			if capacity.PickupDate != "2026-06-11" {
				t.Fatalf("PickupDate = %q, want 2026-06-11", capacity.PickupDate)
			}
			if capacity.MaxOrders != 8 || capacity.CurrentOrders != 3 || capacity.IsClosed {
				t.Fatalf("unexpected capacity: %+v", capacity)
			}
			capacity.ID = 90
			return capacity, nil
		},
	}
	service := pickupcalendar_service.NewPickupCalendarService(repo)

	capacity, err := service.SetCapacity(context.Background(), 2, core_domain.RoleLogist.String(), pickupcalendar_service.SetCapacityInput{
		WarehouseID:   10,
		PickupDate:    "2026-06-11",
		MaxOrders:     8,
		CurrentOrders: 3,
		IsClosed:      false,
	})
	if err != nil {
		t.Fatalf("SetCapacity() error = %v", err)
	}
	if capacity.ID != 90 {
		t.Fatalf("capacity.ID = %d, want 90", capacity.ID)
	}
}

func TestSetCapacityRejectsCurrentGreaterThanMax(t *testing.T) {
	service := pickupcalendar_service.NewPickupCalendarService(&fakePickupCalendarRepository{})
	_, err := service.SetCapacity(context.Background(), 1, core_domain.RoleAdmin.String(), pickupcalendar_service.SetCapacityInput{
		WarehouseID:   10,
		PickupDate:    "2026-06-11",
		MaxOrders:     2,
		CurrentOrders: 3,
	})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("error = %v, want invalid argument", err)
	}
}

func TestGetCalendarRejectsLongRange(t *testing.T) {
	service := pickupcalendar_service.NewPickupCalendarService(&fakePickupCalendarRepository{})
	_, err := service.GetCalendar(context.Background(), 1, core_domain.RoleAdmin.String(), pickupcalendar_service.CalendarFilter{
		WarehouseID: 10,
		DateFrom:    "2026-01-01",
		DateTo:      "2026-12-31",
	})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("error = %v, want invalid argument", err)
	}
}
