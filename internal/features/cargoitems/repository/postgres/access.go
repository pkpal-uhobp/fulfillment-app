package cargoitems_repository_postgres

import (
	"context"
	"fmt"
)

func (r *CargoItemsRepository) ClientOwnsCargoItem(
	ctx context.Context,
	cargoItemID int64,
	clientID int64,
) (bool, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM cargo_items ci
			JOIN orders o ON o.id = ci.order_id
			WHERE ci.id = $1 AND o.client_id = $2
		);
	`

	var exists bool
	if err := q.QueryRow(ctx, query, cargoItemID, clientID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check cargo item ownership: %w", err)
	}

	return exists, nil
}
