package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

func runMigrations(dialect goose.Dialect, db *sql.DB, fs fs.FS) error {
	provider, err := goose.NewProvider(dialect, db, fs)
	if err != nil {
		return fmt.Errorf("Failed to create migration provider: %w", err)
	}
	_, err = provider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}
	return nil
}
