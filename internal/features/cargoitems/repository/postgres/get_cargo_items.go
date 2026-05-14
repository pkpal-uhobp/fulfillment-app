package cargoitems_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

const (
	defaultRepositoryCargoItemsPage  = 1
	defaultRepositoryCargoItemsLimit = 20
)

func (r *CargoItemsRepository) ListCargoItems(
	ctx context.Context,
	filter core_domain.CargoItemFilter,
) ([]core_domain.CargoItem, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	page := filter.Page
	if page <= 0 {
		page = defaultRepositoryCargoItemsPage
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = defaultRepositoryCargoItemsLimit
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT %s
		FROM cargo_items ci
		JOIN orders o ON o.id = ci.order_id
		WHERE ($1::bigint IS NULL OR ci.order_id = $1)
		  AND ($2::varchar = '' OR ci.status = $2)
		  AND ($3::bigint IS NULL OR ci.storage_zone_id = $3)
		  AND ($4::bigint IS NULL OR ci.gate_id = $4)
		  AND ($5::varchar = '' OR ci.qr_code = $5)
		  AND ($6::bigint IS NULL OR o.client_id = $6)
		ORDER BY ci.id DESC
		LIMIT $7 OFFSET $8;
	`, prefixedCargoItemColumns("ci"))

	rows, err := q.Query(
		ctx,
		query,
		filter.OrderID,
		filter.Status,
		filter.StorageZoneID,
		filter.GateID,
		filter.QRCode,
		filter.ClientID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query cargo items: %w", err)
	}
	defer rows.Close()

	items := make([]core_domain.CargoItem, 0)
	for rows.Next() {
		item, err := scanCargoItemRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan cargo item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cargo items: %w", err)
	}

	return items, nil
}

func prefixedCargoItemColumns(prefix string) string {
	return fmt.Sprintf(`
		%s.id,
		%s.order_id,
		%s.order_cargo_place_id,
		%s.cargo_place_type_id,
		%s.qr_code,
		%s.status,
		%s.storage_zone_id,
		%s.gate_id,
		%s.received_by,
		%s.shipped_by,
		%s.received_at,
		%s.shipped_at,
		%s.comment,
		%s.created_at,
		%s.updated_at
	`, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix)
}
