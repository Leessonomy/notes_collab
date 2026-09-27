package app

import (
	"context"
	dbconfig "notes-collab-api/config"
	"notes-collab-api/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var postgresModule = fx.Module("postgres",
	fx.Provide(newPool),
)

func newPool(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func runMigrations(cfg *config.Config) error {
	return dbconfig.Migrate(context.Background(), cfg.DatabaseURL)
}
