package warehouses_repository_postgres

import (
	"context"
	"database/sql"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *WarehousesRepository) ListCargoPlaceTypes(
	ctx context.Context,
) ([]core_domain.CargoPlaceType, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT id, name, description, is_active
		FROM cargo_place_types
		WHERE is_active = true
		ORDER BY name;
	`

	rows, err := q.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list cargo place types: %w", err)
	}
	defer rows.Close()

	cargoPlaceTypes := make([]core_domain.CargoPlaceType, 0)
	for rows.Next() {
		var (
			cargoPlaceType core_domain.CargoPlaceType
			description    sql.NullString
		)

		if err := rows.Scan(
			&cargoPlaceType.ID,
			&cargoPlaceType.Name,
			&description,
			&cargoPlaceType.IsActive,
		); err != nil {
			return nil, fmt.Errorf("scan cargo place type: %w", err)
		}

		if description.Valid {
			cargoPlaceType.Description = &description.String
		}

		cargoPlaceTypes = append(cargoPlaceTypes, cargoPlaceType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cargo place types: %w", err)
	}

	return cargoPlaceTypes, nil
}
