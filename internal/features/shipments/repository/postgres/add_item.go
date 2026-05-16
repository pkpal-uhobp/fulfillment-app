package shipments_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *ShipmentsRepository) AddShipmentItem(
	ctx context.Context,
	shipmentID int64,
	cargoItemID int64,
	changedBy int64,
	comment *string,
) (core_domain.ShipmentDetails, error) {
	var result core_domain.ShipmentDetails
	if err := r.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		q := r.tx.Querier(ctx)

		const shipmentQuery = `
			SELECT destination_warehouse_id, gate_id, status
			FROM shipments
			WHERE id = $1
			FOR UPDATE;
		`
		var shipmentDestinationWarehouseID int64
		var shipmentGateID int64
		var shipmentStatus string
		if err := q.QueryRow(ctx, shipmentQuery, shipmentID).Scan(
			&shipmentDestinationWarehouseID,
			&shipmentGateID,
			&shipmentStatus,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: shipment not found", core_errors.ErrNotFound)
			}
			return fmt.Errorf("select shipment for add item: %w", err)
		}
		if shipmentStatus != core_domain.ShipmentStatusPlanned.String() && shipmentStatus != core_domain.ShipmentStatusLoading.String() {
			return fmt.Errorf("%w: shipment does not accept items in status %s", core_errors.ErrConflict, shipmentStatus)
		}

		const cargoQuery = `
			SELECT ci.status, ci.gate_id, o.destination_warehouse_id
			FROM cargo_items ci
			JOIN orders o ON o.id = ci.order_id
			WHERE ci.id = $1
			FOR UPDATE;
		`
		var cargoStatus string
		var cargoGateID *int64
		var orderDestinationWarehouseID int64
		if err := q.QueryRow(ctx, cargoQuery, cargoItemID).Scan(
			&cargoStatus,
			&cargoGateID,
			&orderDestinationWarehouseID,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: cargo item not found", core_errors.ErrNotFound)
			}
			return fmt.Errorf("select cargo item for shipment: %w", err)
		}
		if orderDestinationWarehouseID != shipmentDestinationWarehouseID {
			return fmt.Errorf("%w: cargo destination warehouse does not match shipment", core_errors.ErrConflict)
		}
		if cargoGateID == nil || *cargoGateID != shipmentGateID {
			return fmt.Errorf("%w: cargo gate does not match shipment gate", core_errors.ErrConflict)
		}
		if cargoStatus != core_domain.CargoItemStatusStored.String() && cargoStatus != core_domain.CargoItemStatusReadyToShip.String() {
			return fmt.Errorf("%w: cargo item must be stored or ready_to_ship", core_errors.ErrConflict)
		}

		const insertQuery = `
			INSERT INTO shipment_items (shipment_id, cargo_item_id)
			VALUES ($1, $2);
		`
		if _, err := q.Exec(ctx, insertQuery, shipmentID, cargoItemID); err != nil {
			if isUniqueViolation(err) {
				return fmt.Errorf("%w: cargo item already added to shipment", core_errors.ErrConflict)
			}
			if isForeignKeyViolation(err) {
				return fmt.Errorf("%w: invalid shipment item foreign key", core_errors.ErrInvalidArgument)
			}
			return fmt.Errorf("insert shipment item: %w", err)
		}

		if cargoStatus != core_domain.CargoItemStatusReadyToShip.String() {
			const updateCargoQuery = `
				UPDATE cargo_items
				SET status = $2,
					updated_at = CURRENT_TIMESTAMP
				WHERE id = $1;
			`
			if _, err := q.Exec(ctx, updateCargoQuery, cargoItemID, core_domain.CargoItemStatusReadyToShip.String()); err != nil {
				return fmt.Errorf("mark cargo item ready to ship: %w", err)
			}
			finalComment := comment
			if finalComment == nil {
				value := "added to shipment"
				finalComment = &value
			}
			const historyQuery = `
				INSERT INTO cargo_status_history (
					cargo_item_id,
					old_status,
					new_status,
					changed_by,
					comment
				) VALUES ($1, $2, $3, $4, $5);
			`
			if _, err := q.Exec(
				ctx,
				historyQuery,
				cargoItemID,
				cargoStatus,
				core_domain.CargoItemStatusReadyToShip.String(),
				changedBy,
				finalComment,
			); err != nil {
				return fmt.Errorf("create cargo ready history: %w", err)
			}
		}

		details, err := getShipmentDetailsByID(ctx, q, shipmentID)
		if err != nil {
			return err
		}
		result = details
		return nil
	}); err != nil {
		return core_domain.ShipmentDetails{}, err
	}
	return result, nil
}
