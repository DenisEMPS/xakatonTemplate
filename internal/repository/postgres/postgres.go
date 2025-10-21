package postgres

import (
	"context"
	"xaxaton/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg config.Postgres) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnTimeout)
	defer cancel()

	connectedPool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := connectedPool.Ping(ctx); err != nil {
		return nil, err
	}

	return connectedPool, nil
}
