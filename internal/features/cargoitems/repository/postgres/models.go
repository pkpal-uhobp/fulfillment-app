package cargoitems_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
)

const cargoItemColumns = `
	id,
	order_id,
	order_cargo_place_id,
	cargo_place_type_id,
	qr_code,
	status,
	storage_zone_id,
	gate_id,
	received_by,
	shipped_by,
	received_at,
	shipped_at,
	comment,
	created_at,
	updated_at
`

func scanCargoItemRow(row pgx.Row) (core_domain.CargoItem, error) {
	var item core_domain.CargoItem
	var status string

	if err := row.Scan(
		&item.ID,
		&item.OrderID,
		&item.OrderCargoPlaceID,
		&item.CargoPlaceTypeID,
		&item.QRCode,
		&status,
		&item.StorageZoneID,
		&item.GateID,
		&item.ReceivedBy,
		&item.ShippedBy,
		&item.ReceivedAt,
		&item.ShippedAt,
		&item.Comment,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return core_domain.CargoItem{}, err
	}

	item.Status = core_domain.CargoItemStatus(status)
	return item, nil
}

func getCargoItemByID(
	ctx context.Context,
	q core_postgres_tx.Querier,
	cargoItemID int64,
) (core_domain.CargoItem, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM cargo_items
		WHERE id = $1;
	`, cargoItemColumns)

	item, err := scanCargoItemRow(q.QueryRow(ctx, query, cargoItemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.CargoItem{}, fmt.Errorf("%w: cargo item not found", core_errors.ErrNotFound)
		}
		return core_domain.CargoItem{}, fmt.Errorf("scan cargo item: %w", err)
	}

	return item, nil
}
