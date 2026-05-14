package cargoitems_transport_http_tests

import (
	"context"
	"fmt"

	cargoitems_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/service"
)

type fakeCargoItemsService struct {
	createCargoItemFn       func(context.Context, int64, string, cargoitems_service.CreateCargoItemInput) (cargoitems_service.CargoItemDTO, error)
	listCargoItemsFn        func(context.Context, int64, string, cargoitems_service.CargoItemFilter) ([]cargoitems_service.CargoItemDTO, error)
	getCargoItemFn          func(context.Context, int64, int64, string) (cargoitems_service.CargoItemDTO, error)
	getCargoItemHistoryFn   func(context.Context, int64, int64, string) ([]cargoitems_service.CargoStatusHistoryDTO, error)
	updateCargoItemStatusFn func(context.Context, int64, int64, string, cargoitems_service.UpdateCargoItemStatusInput) (cargoitems_service.CargoItemDTO, error)
	assignStorageZoneFn     func(context.Context, int64, int64, string, cargoitems_service.AssignStorageZoneInput) (cargoitems_service.CargoItemDTO, error)
	assignGateFn            func(context.Context, int64, int64, string, cargoitems_service.AssignGateInput) (cargoitems_service.CargoItemDTO, error)
}

func (f *fakeCargoItemsService) CreateCargoItem(ctx context.Context, actorID int64, actorRole string, input cargoitems_service.CreateCargoItemInput) (cargoitems_service.CargoItemDTO, error) {
	if f.createCargoItemFn == nil {
		return cargoitems_service.CargoItemDTO{}, fmt.Errorf("unexpected CreateCargoItem call")
	}
	return f.createCargoItemFn(ctx, actorID, actorRole, input)
}

func (f *fakeCargoItemsService) ListCargoItems(ctx context.Context, actorID int64, actorRole string, filter cargoitems_service.CargoItemFilter) ([]cargoitems_service.CargoItemDTO, error) {
	if f.listCargoItemsFn == nil {
		return nil, fmt.Errorf("unexpected ListCargoItems call")
	}
	return f.listCargoItemsFn(ctx, actorID, actorRole, filter)
}

func (f *fakeCargoItemsService) GetCargoItem(ctx context.Context, cargoItemID int64, actorID int64, actorRole string) (cargoitems_service.CargoItemDTO, error) {
	if f.getCargoItemFn == nil {
		return cargoitems_service.CargoItemDTO{}, fmt.Errorf("unexpected GetCargoItem call")
	}
	return f.getCargoItemFn(ctx, cargoItemID, actorID, actorRole)
}

func (f *fakeCargoItemsService) GetCargoItemHistory(ctx context.Context, cargoItemID int64, actorID int64, actorRole string) ([]cargoitems_service.CargoStatusHistoryDTO, error) {
	if f.getCargoItemHistoryFn == nil {
		return nil, fmt.Errorf("unexpected GetCargoItemHistory call")
	}
	return f.getCargoItemHistoryFn(ctx, cargoItemID, actorID, actorRole)
}

func (f *fakeCargoItemsService) UpdateCargoItemStatus(ctx context.Context, cargoItemID int64, actorID int64, actorRole string, input cargoitems_service.UpdateCargoItemStatusInput) (cargoitems_service.CargoItemDTO, error) {
	if f.updateCargoItemStatusFn == nil {
		return cargoitems_service.CargoItemDTO{}, fmt.Errorf("unexpected UpdateCargoItemStatus call")
	}
	return f.updateCargoItemStatusFn(ctx, cargoItemID, actorID, actorRole, input)
}

func (f *fakeCargoItemsService) AssignStorageZone(ctx context.Context, cargoItemID int64, actorID int64, actorRole string, input cargoitems_service.AssignStorageZoneInput) (cargoitems_service.CargoItemDTO, error) {
	if f.assignStorageZoneFn == nil {
		return cargoitems_service.CargoItemDTO{}, fmt.Errorf("unexpected AssignStorageZone call")
	}
	return f.assignStorageZoneFn(ctx, cargoItemID, actorID, actorRole, input)
}

func (f *fakeCargoItemsService) AssignGate(ctx context.Context, cargoItemID int64, actorID int64, actorRole string, input cargoitems_service.AssignGateInput) (cargoitems_service.CargoItemDTO, error) {
	if f.assignGateFn == nil {
		return cargoitems_service.CargoItemDTO{}, fmt.Errorf("unexpected AssignGate call")
	}
	return f.assignGateFn(ctx, cargoItemID, actorID, actorRole, input)
}

func ptr[T any](value T) *T { return &value }
