package core_postgres_tx

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	pool "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/pool"
)

type Tx struct {
	pool      *pool.ConnectionPool
	opTimeout time.Duration
}

func NewTx(db *pool.ConnectionPool) *Tx {
	return &Tx{
		pool:      db,
		opTimeout: db.QueryTimeout(),
	}
}

func (tx *Tx) Querier(ctx context.Context) Querier {
	if q, ok := FromContext(ctx); ok {
		return q
	}

	return tx.pool.Pool
}

func (tx *Tx) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := FromContext(ctx); ok {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, tx.opTimeout)
}

func (tx *Tx) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return tx.WithinTransactionOptions(ctx, pgx.TxOptions{}, fn)
}

func (tx *Tx) WithinTransactionOptions(
	ctx context.Context,
	opts pgx.TxOptions,
	fn func(ctx context.Context) error,
) error {
	if _, ok := FromContext(ctx); ok {
		return fn(ctx)
	}

	txCtx, cancel := context.WithTimeout(ctx, tx.opTimeout)
	defer cancel()

	dbTx, err := tx.pool.BeginTx(txCtx, opts)
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
