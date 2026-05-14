package cargoitems_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (r *CargoItemsRepository) GateBelongsToCargoOrder(
	ctx context.Context,
	cargoItemID int64,
	gateID int64,
) (bool, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM cargo_items ci
			JOIN orders o ON o.id = ci.order_id
			JOIN gates g ON g.id = $2
			WHERE ci.id = $1
			  AND g.warehouse_id = o.receiving_warehouse_id
			  AND g.is_active = TRUE
		);
	`

	var exists bool
	if err := q.QueryRow(ctx, query, cargoItemID, gateID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check gate for cargo item: %w", err)
	}

	return exists, nil
}

func (r *CargoItemsRepository) AssignGate(
	ctx context.Context,
	cargoItemID int64,
	gateID int64,
	changedBy int64,
	comment *string,
) (core_domain.CargoItem, error) {
	var result core_domain.CargoItem

	if err := r.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		q := r.tx.Querier(ctx)

		query := fmt.Sprintf(`
			UPDATE cargo_items
			SET
				gate_id = $2,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
			RETURNING %s;
		`, cargoItemColumns)

		updated, err := scanCargoItemRow(q.QueryRow(ctx, query, cargoItemID, gateID))
		if err != nil {
			if isForeignKeyViolation(err) {
				return fmt.Errorf("%w: invalid gate", core_errors.ErrInvalidArgument)
			}
			return fmt.Errorf("assign gate: %w", err)
		}

		const historyCommentPrefix = "gate assigned"
		finalComment := comment
		if finalComment == nil {
			value := historyCommentPrefix
			finalComment = &value
		}

		const historyQuery = `
			INSERT INTO cargo_status_history (
				cargo_item_id,
				old_status,
				new_status,
				changed_by,
				comment
			)
			VALUES ($1, $2, $2, $3, $4);
		`
		if _, err := q.Exec(
			ctx,
			historyQuery,
			cargoItemID,
			updated.Status.String(),
			changedBy,
			finalComment,
		); err != nil {
			return fmt.Errorf("create cargo assignment history: %w", err)
		}

		result = updated
		return nil
	}); err != nil {
		return core_domain.CargoItem{}, err
	}

	return result, nil
}
