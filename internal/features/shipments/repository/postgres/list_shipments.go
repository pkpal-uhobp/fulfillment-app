package shipments_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

const (
	defaultRepositoryShipmentsPage  = 1
	defaultRepositoryShipmentsLimit = 20
)

func (r *ShipmentsRepository) ListShipments(
	ctx context.Context,
	filter core_domain.ShipmentFilter,
) ([]core_domain.Shipment, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()
	q := r.tx.Querier(ctx)
	page := filter.Page
	if page <= 0 {
		page = defaultRepositoryShipmentsPage
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = defaultRepositoryShipmentsLimit
	}
	offset := (page - 1) * limit
	query := fmt.Sprintf(`
		SELECT %s
		FROM shipments s
		WHERE ($1::varchar = '' OR s.status = $1)
			AND ($2::bigint IS NULL OR s.destination_warehouse_id = $2)
			AND ($3::bigint IS NULL OR s.gate_id = $3)
			AND ($4::date IS NULL OR s.planned_departure_at::date = $4::date)
		ORDER BY s.id DESC
		LIMIT $5 OFFSET $6;
	`, prefixedShipmentColumns("s"))
	var dateArg any
	if filter.Date != "" {
		dateArg = filter.Date
	}
	rows, err := q.Query(
		ctx,
		query,
		filter.Status,
		filter.DestinationWarehouseID,
		filter.GateID,
		dateArg,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query shipments: %w", err)
	}
	defer rows.Close()
	shipments := make([]core_domain.Shipment, 0)
	for rows.Next() {
		shipment, err := scanShipmentRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan shipment: %w", err)
		}
		shipments = append(shipments, shipment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shipments: %w", err)
	}
	return shipments, nil
}
