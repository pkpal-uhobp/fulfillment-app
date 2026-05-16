package shipments_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *ShipmentsRepository) RemoveShipmentItem(
	ctx context.Context,
	shipmentID int64,
	cargoItemID int64,
) error {
	return r.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		q := r.tx.Querier(ctx)
		const shipmentQuery = `
			SELECT status
			FROM shipments
			WHERE id = $1
			FOR UPDATE;
		`
		var status string
		if err := q.QueryRow(ctx, shipmentQuery, shipmentID).Scan(&status); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: shipment not found", core_errors.ErrNotFound)
			}
			return fmt.Errorf("select shipment: %w", err)
		}
		if status != core_domain.ShipmentStatusPlanned.String() && status != core_domain.ShipmentStatusLoading.String() {
			return fmt.Errorf("%w: shipment item cannot be removed from status %s", core_errors.ErrConflict, status)
		}
		const deleteQuery = `
			DELETE FROM shipment_items
			WHERE shipment_id = $1 AND cargo_item_id = $2;
		`
		result, err := q.Exec(ctx, deleteQuery, shipmentID, cargoItemID)
		if err != nil {
			return fmt.Errorf("delete shipment item: %w", err)
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("%w: shipment item not found", core_errors.ErrNotFound)
		}
		return nil
	})
}
