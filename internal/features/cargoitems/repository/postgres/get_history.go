package cargoitems_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *CargoItemsRepository) ListCargoStatusHistory(
	ctx context.Context,
	cargoItemID int64,
) ([]core_domain.CargoStatusHistory, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT
			id,
			cargo_item_id,
			old_status,
			new_status,
			changed_by,
			comment,
			changed_at
		FROM cargo_status_history
		WHERE cargo_item_id = $1
		ORDER BY changed_at ASC, id ASC;
	`

	rows, err := q.Query(ctx, query, cargoItemID)
	if err != nil {
		return nil, fmt.Errorf("query cargo status history: %w", err)
	}
	defer rows.Close()

	history := make([]core_domain.CargoStatusHistory, 0)
	for rows.Next() {
		var item core_domain.CargoStatusHistory
		var oldStatus *string
		var newStatus string

		if err := rows.Scan(
			&item.ID,
			&item.CargoItemID,
			&oldStatus,
			&newStatus,
			&item.ChangedBy,
			&item.Comment,
			&item.ChangedAt,
		); err != nil {
			return nil, fmt.Errorf("scan cargo status history: %w", err)
		}

		if oldStatus != nil {
			value := core_domain.CargoItemStatus(*oldStatus)
			item.OldStatus = &value
		}
		item.NewStatus = core_domain.CargoItemStatus(newStatus)

		history = append(history, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cargo status history: %w", err)
	}

	return history, nil
}
