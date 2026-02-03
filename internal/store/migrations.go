package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sync"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

var gooseOnce sync.Once
var errGooseInit error

func initGoose() error {
	gooseOnce.Do(func() {
		goose.SetBaseFS(migrations)
		errGooseInit = goose.SetDialect("sqlite")
	})
	return errGooseInit
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	if err := initGoose(); err != nil {
		return fmt.Errorf("init goose: %w", err)
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
