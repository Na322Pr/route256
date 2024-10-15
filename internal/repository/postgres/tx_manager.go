package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txManagerKey struct{}

type TxManager struct {
	pool *pgxpool.Pool
}

type txFn = func(ctxTx context.Context) error

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (m *TxManager) RunSerializable(ctx context.Context, fn txFn) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadOnly,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) RunReadCommitted(ctx context.Context, fn txFn) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) beginFunc(ctx context.Context, opts pgx.TxOptions, fn txFn) error {
	tx, err := m.pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() error {
		err = errors.Join(err, tx.Rollback(ctx))
		if err != nil {
			return err
		}
		return nil
	}()

	ctx = context.WithValue(ctx, txManagerKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (m *TxManager) GetQueryEngine(ctx context.Context) QueryEngine {
	v, ok := ctx.Value(txManagerKey{}).(QueryEngine)
	if ok && v != nil {
		return v
	}

	return m.pool
}
