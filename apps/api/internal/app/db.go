package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createPostgres(ctx context.Context, url string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, url)
}
