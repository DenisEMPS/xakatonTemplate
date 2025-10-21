package repository

import (
	"context"

	"xaxaton/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xakepp35/pkg/xlog"
)

type TxOperation func(repo Repository) error

type Repository interface {
	sqlc.Querier
}

type repository struct {
	*sqlc.Queries
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repository {
	return &repository{
		Queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *repository) ExecWithTx(ctx context.Context, fn TxOperation) (err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				xlog.ErrWarn(rollbackErr).Msg("rollback failed")
			}
		}
	}()

	txRepo := &repository{
		Queries: r.Queries.WithTx(tx),
		pool:    r.pool,
	}

	err = fn(txRepo)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)

	return nil
}
