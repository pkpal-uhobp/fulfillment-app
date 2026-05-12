package core_postgres_tx

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Tx struct {
	pool         *pool.ConnectionPool
	queryTimeout time.Duration
}

func NewTx(db *pool.ConnectionPool) *Tx {
	return &Tx{
		pool:         db,
		queryTimeout: db.QueryTimeout(),
	}
}

func (m *Tx) Querier(ctx context.Context) Querier {
	if q, ok := FromContext(ctx); ok {
		return q
	}

	return m.pool.Pool
}

func (m *Tx) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return m.WithinTransactionOptions(ctx, pgx.TxOptions{}, fn)
}

func (m *Tx) WithinTransactionOptions(
	ctx context.Context,
	opts pgx.TxOptions,
	fn func(ctx context.Context) error,
) error {
	if _, ok := FromContext(ctx); ok {
		return fn(ctx)
	}

	txCtx, cancel := context.WithTimeout(ctx, m.queryTimeout)
	defer cancel()

	dbTx, err := m.pool.BeginTx(txCtx, opts)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx = injectTx(txCtx, dbTx)

	defer func() {
		_ = dbTx.Rollback(context.Background())
	}()

	if err := fn(txCtx); err != nil {
		return err
	}

	if err := dbTx.Commit(txCtx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
