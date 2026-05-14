package warehouses_repository_postgres

import (
	"context"
	"database/sql"
	"fmt"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (r *WarehousesRepository) ListProductTypes(
	ctx context.Context,
) ([]core_domain.ProductType, error) {
	ctx, cancel := r.tx.WithTimeout(ctx)
	defer cancel()

	q := r.tx.Querier(ctx)

	const query = `
		SELECT id, name, description, is_active
		FROM product_types
		WHERE is_active = true
		ORDER BY name;
	`

	rows, err := q.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list product types: %w", err)
	}
	defer rows.Close()

	productTypes := make([]core_domain.ProductType, 0)
	for rows.Next() {
		var (
			productType core_domain.ProductType
			description sql.NullString
		)

		if err := rows.Scan(
			&productType.ID,
			&productType.Name,
			&description,
			&productType.IsActive,
		); err != nil {
			return nil, fmt.Errorf("scan product type: %w", err)
		}

		if description.Valid {
			productType.Description = &description.String
		}

		productTypes = append(productTypes, productType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product types: %w", err)
	}

	return productTypes, nil
}
