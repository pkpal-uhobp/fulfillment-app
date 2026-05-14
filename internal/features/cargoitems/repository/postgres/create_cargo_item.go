package cargoitems_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *CargoItemsRepository) CreateCargoItem(
	ctx context.Context,
	item core_domain.CargoItem,
	changedBy int64,
) (core_domain.CargoItem, error) {
	var result core_domain.CargoItem

	if err := r.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		q := r.tx.Querier(ctx)

		query := fmt.Sprintf(`
			INSERT INTO cargo_items (
				order_id,
				order_cargo_place_id,
				cargo_place_type_id,
				qr_code,
				status,
				received_by,
				received_at,
				comment
			)
			VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, $7)
			RETURNING %s;
		`, cargoItemColumns)

		created, err := scanCargoItemRow(q.QueryRow(
			ctx,
			query,
			item.OrderID,
			item.OrderCargoPlaceID,
			item.CargoPlaceTypeID,
			item.QRCode,
			item.Status.String(),
			changedBy,
			item.Comment,
		))
		if err != nil {
			switch {
			case isUniqueViolation(err):
				return fmt.Errorf("%w: cargo item qr code already exists", core_errors.ErrConflict)
			case isForeignKeyViolation(err), isCheckViolation(err):
				return fmt.Errorf("%w: invalid cargo item data", core_errors.ErrInvalidArgument)
			default:
				return fmt.Errorf("create cargo item: %w", err)
			}
		}

		const historyQuery = `
			INSERT INTO cargo_status_history (
				cargo_item_id,
				old_status,
				new_status,
				changed_by,
				comment
			)
			VALUES ($1, NULL, $2, $3, $4);
		`

		if _, err := q.Exec(
			ctx,
			historyQuery,
			created.ID,
			created.Status.String(),
			changedBy,
			created.Comment,
		); err != nil {
			return fmt.Errorf("create cargo status history: %w", err)
		}

		result = created
		return nil
	}); err != nil {
		return core_domain.CargoItem{}, err
	}

	return result, nil
}
