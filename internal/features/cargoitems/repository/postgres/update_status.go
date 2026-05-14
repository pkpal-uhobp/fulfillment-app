package cargoitems_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	coreDomain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	coreErrors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *CargoItemsRepository) UpdateCargoItemStatus(
	ctx context.Context,
	cargoItemID int64,
	status string,
	changedBy int64,
	comment *string,
) (coreDomain.CargoItem, error) {
	var result coreDomain.CargoItem

	if err := r.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		q := r.tx.Querier(ctx)

		const selectQuery = `
			SELECT status
			FROM cargo_items
			WHERE id = $1
			FOR UPDATE;
		`

		var oldStatus string
		if err := q.QueryRow(ctx, selectQuery, cargoItemID).Scan(&oldStatus); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: cargo item not found", coreErrors.ErrNotFound)
			}
			return fmt.Errorf("select cargo item status: %w", err)
		}

		updateQuery := fmt.Sprintf(`
			UPDATE cargo_items
			SET
				status = $2::varchar(50),
				shipped_by = CASE
					WHEN $2::varchar(50) = 'shipped' THEN $3
					ELSE shipped_by
				END,
				shipped_at = CASE
					WHEN $2::varchar(50) = 'shipped' THEN CURRENT_TIMESTAMP
					ELSE shipped_at
				END,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
			RETURNING %s;
		`, cargoItemColumns)

		updated, err := scanCargoItemRow(q.QueryRow(
			ctx,
			updateQuery,
			cargoItemID,
			status,
			changedBy,
		))
		if err != nil {
			if isCheckViolation(err) {
				return fmt.Errorf("%w: invalid cargo item status", coreErrors.ErrInvalidArgument)
			}
			return fmt.Errorf("update cargo item status: %w", err)
		}

		const historyQuery = `
			INSERT INTO cargo_status_history (
				cargo_item_id,
				old_status,
				new_status,
				changed_by,
				comment
			)
			VALUES ($1, $2, $3, $4, $5);
		`

		if _, err := q.Exec(
			ctx,
			historyQuery,
			cargoItemID,
			oldStatus,
			status,
			changedBy,
			comment,
		); err != nil {
			return fmt.Errorf("create cargo status history: %w", err)
		}

		result = updated
		return nil
	}); err != nil {
		return coreDomain.CargoItem{}, err
	}

	return result, nil
}
