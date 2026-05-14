package cargoitems_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *CargoItemsRepository) GetCargoItem(
	ctx context.Context,
	cargoItemID int64,
) (core_domain.CargoItem, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	item, err := getCargoItemByID(ctx, q, cargoItemID)
	if err != nil {
		return core_domain.CargoItem{}, fmt.Errorf("get cargo item: %w", err)
	}

	return item, nil
}
