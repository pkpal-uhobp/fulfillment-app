package shipments_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *ShipmentsRepository) CreateShipment(
	ctx context.Context,
	shipment core_domain.Shipment,
) (core_domain.Shipment, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()
	q := r.tx.Querier(ctx)
	query := fmt.Sprintf(`
		INSERT INTO shipments (
			destination_warehouse_id,
			gate_id,
			planned_departure_at,
			status,
			created_by
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING %s;
	`, shipmentColumns)
	created, err := scanShipmentRow(q.QueryRow(
		ctx,
		query,
		shipment.DestinationWarehouseID,
		shipment.GateID,
		shipment.PlannedDepartureAt,
		shipment.Status.String(),
		shipment.CreatedBy,
	))
	if err != nil {
		if isForeignKeyViolation(err) {
			return core_domain.Shipment{}, fmt.Errorf("%w: invalid shipment foreign key", core_errors.ErrInvalidArgument)
		}
		if isCheckViolation(err) {
			return core_domain.Shipment{}, fmt.Errorf("%w: invalid shipment data", core_errors.ErrInvalidArgument)
		}
		return core_domain.Shipment{}, fmt.Errorf("create shipment: %w", err)
	}
	return created, nil
}
