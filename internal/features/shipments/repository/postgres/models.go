package shipments_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
)

const shipmentColumns = `
	id,
	destination_warehouse_id,
	gate_id,
	planned_departure_at,
	actual_departure_at,
	status,
	created_by,
	created_at
`

const shipmentItemColumns = `
	si.id,
	si.shipment_id,
	si.cargo_item_id,
	ci.id,
	ci.order_id,
	ci.order_cargo_place_id,
	ci.cargo_place_type_id,
	ci.qr_code,
	ci.status,
	ci.storage_zone_id,
	ci.gate_id,
	ci.received_by,
	ci.shipped_by,
	ci.received_at,
	ci.shipped_at,
	ci.comment,
	ci.created_at,
	ci.updated_at
`

func prefixedShipmentColumns(prefix string) string {
	return fmt.Sprintf(`
		%s.id,
		%s.destination_warehouse_id,
		%s.gate_id,
		%s.planned_departure_at,
		%s.actual_departure_at,
		%s.status,
		%s.created_by,
		%s.created_at
	`, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix)
}

func scanShipmentRow(row pgx.Row) (core_domain.Shipment, error) {
	var shipment core_domain.Shipment
	var status string
	if err := row.Scan(
		&shipment.ID,
		&shipment.DestinationWarehouseID,
		&shipment.GateID,
		&shipment.PlannedDepartureAt,
		&shipment.ActualDepartureAt,
		&status,
		&shipment.CreatedBy,
		&shipment.CreatedAt,
	); err != nil {
		return core_domain.Shipment{}, err
	}
	shipment.Status = core_domain.ShipmentStatus(status)
	return shipment, nil
}

func scanShipmentItemDetailsRow(row pgx.Row) (core_domain.ShipmentItemDetails, error) {
	var details core_domain.ShipmentItemDetails
	var status string
	if err := row.Scan(
		&details.Item.ID,
		&details.Item.ShipmentID,
		&details.Item.CargoItemID,
		&details.CargoItem.ID,
		&details.CargoItem.OrderID,
		&details.CargoItem.OrderCargoPlaceID,
		&details.CargoItem.CargoPlaceTypeID,
		&details.CargoItem.QRCode,
		&status,
		&details.CargoItem.StorageZoneID,
		&details.CargoItem.GateID,
		&details.CargoItem.ReceivedBy,
		&details.CargoItem.ShippedBy,
		&details.CargoItem.ReceivedAt,
		&details.CargoItem.ShippedAt,
		&details.CargoItem.Comment,
		&details.CargoItem.CreatedAt,
		&details.CargoItem.UpdatedAt,
	); err != nil {
		return core_domain.ShipmentItemDetails{}, err
	}
	details.CargoItem.Status = core_domain.CargoItemStatus(status)
	return details, nil
}

func getShipmentByID(
	ctx context.Context,
	q core_postgres_tx.Querier,
	shipmentID int64,
) (core_domain.Shipment, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM shipments
		WHERE id = $1;
	`, shipmentColumns)
	shipment, err := scanShipmentRow(q.QueryRow(ctx, query, shipmentID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.Shipment{}, fmt.Errorf("%w: shipment not found", core_errors.ErrNotFound)
		}
		return core_domain.Shipment{}, fmt.Errorf("scan shipment: %w", err)
	}
	return shipment, nil
}

func getShipmentDetailsByID(
	ctx context.Context,
	q core_postgres_tx.Querier,
	shipmentID int64,
) (core_domain.ShipmentDetails, error) {
	shipment, err := getShipmentByID(ctx, q, shipmentID)
	if err != nil {
		return core_domain.ShipmentDetails{}, err
	}
	items, err := listShipmentItems(ctx, q, shipmentID)
	if err != nil {
		return core_domain.ShipmentDetails{}, err
	}
	return core_domain.ShipmentDetails{
		Shipment: shipment,
		Items:    items,
	}, nil
}

func listShipmentItems(
	ctx context.Context,
	q core_postgres_tx.Querier,
	shipmentID int64,
) ([]core_domain.ShipmentItemDetails, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM shipment_items si
		JOIN cargo_items ci ON ci.id = si.cargo_item_id
		WHERE si.shipment_id = $1
		ORDER BY si.id;
	`, shipmentItemColumns)
	rows, err := q.Query(ctx, query, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("query shipment items: %w", err)
	}
	defer rows.Close()
	items := make([]core_domain.ShipmentItemDetails, 0)
	for rows.Next() {
		item, err := scanShipmentItemDetailsRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan shipment item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shipment items: %w", err)
	}
	return items, nil
}
