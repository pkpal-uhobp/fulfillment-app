package shipments_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *ShipmentsRepository) GetShipment(
	ctx context.Context,
	shipmentID int64,
) (core_domain.ShipmentDetails, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()
	q := r.tx.Querier(ctx)
	details, err := getShipmentDetailsByID(ctx, q, shipmentID)
	if err != nil {
		return core_domain.ShipmentDetails{}, fmt.Errorf("get shipment: %w", err)
	}
	return details, nil
}
