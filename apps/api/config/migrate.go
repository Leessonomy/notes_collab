package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"notes-collab-api/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open db for migrations: %w", err)
	}
	defer db.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return fmt.Errorf("create goose provider: %w", err)
	}

	log.Println("migrations: checking for pending migrations")

	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	if len(results) == 0 {
		log.Println("migrations: already up to date, nothing to apply")
		return nil
	}

	for _, r := range results {
		log.Printf("migrations: applied %s (%s)", r.Source.Path, r.Duration)
	}
	log.Printf("migrations: done, applied %d migration(s)", len(results))

	return nil
}
