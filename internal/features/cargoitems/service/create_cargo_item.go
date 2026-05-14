package cargoitems_service

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *CargoItemsService) CreateCargoItem(
	ctx context.Context,
	actorID int64,
	actorRole string,
	input CreateCargoItemInput,
) (CargoItemDTO, error) {
	if actorID <= 0 {
		return CargoItemDTO{}, fmt.Errorf("%w: invalid actor id", core_errors.ErrInvalidArgument)
	}
	if actorRole != core_domain.RoleWorker.String() && actorRole != core_domain.RoleAdmin.String() {
		return CargoItemDTO{}, fmt.Errorf("%w: only worker or admin can create cargo items", core_errors.ErrForbidden)
	}
	if input.OrderID <= 0 {
		return CargoItemDTO{}, fmt.Errorf("%w: order id is required", core_errors.ErrInvalidArgument)
	}
	if input.OrderCargoPlaceID <= 0 {
		return CargoItemDTO{}, fmt.Errorf("%w: order cargo place id is required", core_errors.ErrInvalidArgument)
	}

	cargoPlace, orderStatus, err := s.repo.GetOrderCargoPlaceForOrder(
		ctx,
		input.OrderID,
		input.OrderCargoPlaceID,
	)
	if err != nil {
		return CargoItemDTO{}, err
	}

	if orderStatus == core_domain.OrderStatusCancelled || orderStatus == core_domain.OrderStatusDelivered {
		return CargoItemDTO{}, fmt.Errorf("%w: cannot create cargo item for terminal order", core_errors.ErrConflict)
	}

	createdCount, err := s.repo.CountCargoItemsByOrderCargoPlace(ctx, input.OrderCargoPlaceID)
	if err != nil {
		return CargoItemDTO{}, err
	}
	if createdCount >= cargoPlace.Quantity {
		return CargoItemDTO{}, fmt.Errorf("%w: declared cargo place quantity is already accepted", core_errors.ErrConflict)
	}

	qrCode, err := normalizeQRCode(input.QRCode)
	if err != nil {
		return CargoItemDTO{}, err
	}

	item := core_domain.CargoItem{
		OrderID:           input.OrderID,
		OrderCargoPlaceID: input.OrderCargoPlaceID,
		CargoPlaceTypeID:  cargoPlace.CargoPlaceTypeID,
		QRCode:            qrCode,
		Status:            core_domain.CargoItemStatusAccepted,
		Comment:           normalizeOptionalString(input.Comment),
	}

	created, err := s.repo.CreateCargoItem(ctx, item, actorID)
	if err != nil {
		return CargoItemDTO{}, err
	}

	return mapCargoItemToDTO(created), nil
}
