package orders_repository_postgres

import core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"

type OrdersRepository struct {
	tx *core_postgres_tx.Tx
}

func NewOrdersRepository(tx *core_postgres_tx.Tx) *OrdersRepository {
	return &OrdersRepository{
		tx: tx,
	}
}
