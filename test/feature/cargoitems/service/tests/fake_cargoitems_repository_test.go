package cargoitems_service_tests

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type fakeCargoItemsRepository struct {
	createCargoItemFn                  func(context.Context, core_domain.CargoItem, int64) (core_domain.CargoItem, error)
	listCargoItemsFn                   func(context.Context, core_domain.CargoItemFilter) ([]core_domain.CargoItem, error)
	getCargoItemFn                     func(context.Context, int64) (core_domain.CargoItem, error)
	listCargoStatusHistoryFn           func(context.Context, int64) ([]core_domain.CargoStatusHistory, error)
	updateCargoItemStatusFn            func(context.Context, int64, string, int64, *string) (core_domain.CargoItem, error)
	assignStorageZoneFn                func(context.Context, int64, int64, int64, *string) (core_domain.CargoItem, error)
	assignGateFn                       func(context.Context, int64, int64, int64, *string) (core_domain.CargoItem, error)
	getOrderCargoPlaceForOrderFn       func(context.Context, int64, int64) (core_domain.OrderCargoPlace, core_domain.OrderStatus, error)
	countCargoItemsByOrderCargoPlaceFn func(context.Context, int64) (int, error)
	clientOwnsCargoItemFn              func(context.Context, int64, int64) (bool, error)
	storageZoneBelongsToCargoOrderFn   func(context.Context, int64, int64) (bool, error)
	gateBelongsToCargoOrderFn          func(context.Context, int64, int64) (bool, error)
}

func (f *fakeCargoItemsRepository) CreateCargoItem(ctx context.Context, item core_domain.CargoItem, changedBy int64) (core_domain.CargoItem, error) {
	if f.createCargoItemFn == nil {
		return core_domain.CargoItem{}, fmt.Errorf("unexpected CreateCargoItem call")
	}
	return f.createCargoItemFn(ctx, item, changedBy)
}

func (f *fakeCargoItemsRepository) ListCargoItems(ctx context.Context, filter core_domain.CargoItemFilter) ([]core_domain.CargoItem, error) {
	if f.listCargoItemsFn == nil {
		return nil, fmt.Errorf("unexpected ListCargoItems call")
	}
	return f.listCargoItemsFn(ctx, filter)
}

func (f *fakeCargoItemsRepository) GetCargoItem(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
	if f.getCargoItemFn == nil {
		return core_domain.CargoItem{}, fmt.Errorf("unexpected GetCargoItem call")
	}
	return f.getCargoItemFn(ctx, cargoItemID)
}

func (f *fakeCargoItemsRepository) ListCargoStatusHistory(ctx context.Context, cargoItemID int64) ([]core_domain.CargoStatusHistory, error) {
	if f.listCargoStatusHistoryFn == nil {
		return nil, fmt.Errorf("unexpected ListCargoStatusHistory call")
	}
	return f.listCargoStatusHistoryFn(ctx, cargoItemID)
}

func (f *fakeCargoItemsRepository) UpdateCargoItemStatus(ctx context.Context, cargoItemID int64, status string, changedBy int64, comment *string) (core_domain.CargoItem, error) {
	if f.updateCargoItemStatusFn == nil {
		return core_domain.CargoItem{}, fmt.Errorf("unexpected UpdateCargoItemStatus call")
	}
	return f.updateCargoItemStatusFn(ctx, cargoItemID, status, changedBy, comment)
}

func (f *fakeCargoItemsRepository) AssignStorageZone(ctx context.Context, cargoItemID int64, storageZoneID int64, changedBy int64, comment *string) (core_domain.CargoItem, error) {
	if f.assignStorageZoneFn == nil {
		return core_domain.CargoItem{}, fmt.Errorf("unexpected AssignStorageZone call")
	}
	return f.assignStorageZoneFn(ctx, cargoItemID, storageZoneID, changedBy, comment)
}

func (f *fakeCargoItemsRepository) AssignGate(ctx context.Context, cargoItemID int64, gateID int64, changedBy int64, comment *string) (core_domain.CargoItem, error) {
	if f.assignGateFn == nil {
		return core_domain.CargoItem{}, fmt.Errorf("unexpected AssignGate call")
	}
	return f.assignGateFn(ctx, cargoItemID, gateID, changedBy, comment)
}

func (f *fakeCargoItemsRepository) GetOrderCargoPlaceForOrder(ctx context.Context, orderID int64, orderCargoPlaceID int64) (core_domain.OrderCargoPlace, core_domain.OrderStatus, error) {
	if f.getOrderCargoPlaceForOrderFn == nil {
		return core_domain.OrderCargoPlace{}, "", fmt.Errorf("unexpected GetOrderCargoPlaceForOrder call")
	}
	return f.getOrderCargoPlaceForOrderFn(ctx, orderID, orderCargoPlaceID)
}

func (f *fakeCargoItemsRepository) CountCargoItemsByOrderCargoPlace(ctx context.Context, orderCargoPlaceID int64) (int, error) {
	if f.countCargoItemsByOrderCargoPlaceFn == nil {
		return 0, fmt.Errorf("unexpected CountCargoItemsByOrderCargoPlace call")
	}
	return f.countCargoItemsByOrderCargoPlaceFn(ctx, orderCargoPlaceID)
}

func (f *fakeCargoItemsRepository) ClientOwnsCargoItem(ctx context.Context, cargoItemID int64, clientID int64) (bool, error) {
	if f.clientOwnsCargoItemFn == nil {
		return false, fmt.Errorf("unexpected ClientOwnsCargoItem call")
	}
	return f.clientOwnsCargoItemFn(ctx, cargoItemID, clientID)
}

func (f *fakeCargoItemsRepository) StorageZoneBelongsToCargoOrder(ctx context.Context, cargoItemID int64, storageZoneID int64) (bool, error) {
	if f.storageZoneBelongsToCargoOrderFn == nil {
		return false, fmt.Errorf("unexpected StorageZoneBelongsToCargoOrder call")
	}
	return f.storageZoneBelongsToCargoOrderFn(ctx, cargoItemID, storageZoneID)
}

func (f *fakeCargoItemsRepository) GateBelongsToCargoOrder(ctx context.Context, cargoItemID int64, gateID int64) (bool, error) {
	if f.gateBelongsToCargoOrderFn == nil {
		return false, fmt.Errorf("unexpected GateBelongsToCargoOrder call")
	}
	return f.gateBelongsToCargoOrderFn(ctx, cargoItemID, gateID)
}

func ptr[T any](value T) *T { return &value }
