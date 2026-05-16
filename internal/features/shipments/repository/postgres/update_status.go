package shipments_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *ShipmentsRepository) UpdateShipmentStatus(
	ctx context.Context,
	shipmentID int64,
	status string,
	changedBy int64,
	comment *string,
) (core_domain.ShipmentDetails, error) {
	var result core_domain.ShipmentDetails
	if err := r.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		q := r.tx.Querier(ctx)

		const lockQuery = `
			SELECT status
			FROM shipments
			WHERE id = $1
			FOR UPDATE;
		`
		var oldStatus string
		if err := q.QueryRow(ctx, lockQuery, shipmentID).Scan(&oldStatus); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: shipment not found", core_errors.ErrNotFound)
			}
			return fmt.Errorf("select shipment status: %w", err)
		}

		if status == core_domain.ShipmentStatusShipped.String() {
			const countQuery = `
				SELECT COUNT(*)
				FROM shipment_items
				WHERE shipment_id = $1;
			`
			var itemsCount int
			if err := q.QueryRow(ctx, countQuery, shipmentID).Scan(&itemsCount); err != nil {
				return fmt.Errorf("count shipment items: %w", err)
			}
			if itemsCount == 0 {
				return fmt.Errorf("%w: shipment cannot be shipped without cargo items", core_errors.ErrConflict)
			}
		}

		updateQuery := fmt.Sprintf(`
			UPDATE shipments
			SET status = $2,
				actual_departure_at = CASE WHEN $2 = 'shipped' THEN CURRENT_TIMESTAMP ELSE actual_departure_at END
			WHERE id = $1
			RETURNING %s;
		`, shipmentColumns)
		if _, err := scanShipmentRow(q.QueryRow(ctx, updateQuery, shipmentID, status)); err != nil {
			if isCheckViolation(err) {
				return fmt.Errorf("%w: invalid shipment status", core_errors.ErrInvalidArgument)
			}
			return fmt.Errorf("update shipment status: %w", err)
		}

		if status == core_domain.ShipmentStatusShipped.String() {
			if err := r.markShipmentCargoItemsShipped(ctx, shipmentID, changedBy, comment); err != nil {
				return err
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

func (r *ShipmentsRepository) markShipmentCargoItemsShipped(
	ctx context.Context,
	shipmentID int64,
	changedBy int64,
	comment *string,
) error {
	q := r.tx.Querier(ctx)
	const selectQuery = `
		SELECT ci.id, ci.status
		FROM cargo_items ci
		JOIN shipment_items si ON si.cargo_item_id = ci.id
		WHERE si.shipment_id = $1
			AND ci.status <> 'shipped'
		FOR UPDATE;
	`
	rows, err := q.Query(ctx, selectQuery, shipmentID)
	if err != nil {
		return fmt.Errorf("select shipment cargo items: %w", err)
	}
	defer rows.Close()
	type cargoStatus struct {
		id     int64
		status string
	}
	items := make([]cargoStatus, 0)
	for rows.Next() {
		var item cargoStatus
		if err := rows.Scan(&item.id, &item.status); err != nil {
			return fmt.Errorf("scan shipment cargo item status: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate shipment cargo items: %w", err)
	}

	finalComment := comment
	if finalComment == nil {
		value := "shipment shipped"
		finalComment = &value
	}
	for _, item := range items {
		const updateCargoQuery = `
			UPDATE cargo_items
			SET status = 'shipped',
				shipped_by = $2,
				shipped_at = CURRENT_TIMESTAMP,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1;
		`
		if _, err := q.Exec(ctx, updateCargoQuery, item.id, changedBy); err != nil {
			return fmt.Errorf("mark cargo item shipped: %w", err)
		}
		const historyQuery = `
			INSERT INTO cargo_status_history (
				cargo_item_id,
				old_status,
				new_status,
				changed_by,
				comment
			) VALUES ($1, $2, 'shipped', $3, $4);
		`
		if _, err := q.Exec(ctx, historyQuery, item.id, item.status, changedBy, finalComment); err != nil {
			return fmt.Errorf("create cargo shipped history: %w", err)
		}
	}
	return nil
}
