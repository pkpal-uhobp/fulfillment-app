package cargoitems_service

import (
	"context"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

type CargoItemsRepository interface {
	CreateCargoItem(
		ctx context.Context,
		item core_domain.CargoItem,
		changedBy int64,
	) (core_domain.CargoItem, error)

	ListCargoItems(
		ctx context.Context,
		filter core_domain.CargoItemFilter,
	) ([]core_domain.CargoItem, error)

	GetCargoItem(
		ctx context.Context,
		cargoItemID int64,
	) (core_domain.CargoItem, error)

	ListCargoStatusHistory(
		ctx context.Context,
		cargoItemID int64,
	) ([]core_domain.CargoStatusHistory, error)

	UpdateCargoItemStatus(
		ctx context.Context,
		cargoItemID int64,
		status string,
		changedBy int64,
		comment *string,
	) (core_domain.CargoItem, error)

	AssignStorageZone(
		ctx context.Context,
		cargoItemID int64,
		storageZoneID int64,
		changedBy int64,
		comment *string,
	) (core_domain.CargoItem, error)

	AssignGate(
		ctx context.Context,
		cargoItemID int64,
		gateID int64,
		changedBy int64,
		comment *string,
	) (core_domain.CargoItem, error)

	GetOrderCargoPlaceForOrder(
		ctx context.Context,
		orderID int64,
		orderCargoPlaceID int64,
	) (core_domain.OrderCargoPlace, core_domain.OrderStatus, error)

	CountCargoItemsByOrderCargoPlace(
		ctx context.Context,
		orderCargoPlaceID int64,
	) (int, error)

	ClientOwnsCargoItem(
		ctx context.Context,
		cargoItemID int64,
		clientID int64,
	) (bool, error)

	StorageZoneBelongsToCargoOrder(
		ctx context.Context,
		cargoItemID int64,
		storageZoneID int64,
	) (bool, error)

	GateBelongsToCargoOrder(
		ctx context.Context,
		cargoItemID int64,
		gateID int64,
	) (bool, error)
}

type CargoItemsService struct {
	repo CargoItemsRepository
}

func NewCargoItemsService(repo CargoItemsRepository) *CargoItemsService {
	return &CargoItemsService{repo: repo}
}
