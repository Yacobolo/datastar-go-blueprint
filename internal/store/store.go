// Package store provides SQLite-based state management with transaction support.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // SQLite driver registration

	"github.com/yacobolo/datastar-go-blueprint/internal/store/queries"
)

// txKey is the context key for storing transaction state.
type txKey struct{}

// SQLiteStore wraps the database connection and queries.
type SQLiteStore struct {
	db      *sql.DB
	queries *queries.Queries
}

// Open creates a new SQLiteStore with the given DSN.
// DSN examples: ":memory:", "file:todos.db", "./data/todos.db"
func Open(dsn string) (*SQLiteStore, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dsn)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better concurrency
	ctx := context.Background()
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA busy_timeout = 5000",
	}
	for _, pragma := range pragmas {
		if _, err := db.ExecContext(ctx, pragma); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("failed to set pragma: %w", err)
		}
	}

	// Run migrations
	if err := runMigrations(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &SQLiteStore{
		db:      db,
		queries: queries.New(db),
	}, nil
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Queries returns the sqlc-generated queries.
func (s *SQLiteStore) Queries() *queries.Queries {
	return s.queries
}

// DB returns the underlying database connection for transactions.
func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}

// WithinTransaction executes fn within a database transaction.
// If fn returns an error, the transaction is rolled back.
// If fn returns nil, the transaction is committed.
// The txCtx carries the transaction state for repositories to use.
func (s *SQLiteStore) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Embed transaction in context
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Join(err, fmt.Errorf("rollback failed: %w", rbErr))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// conn returns the appropriate queries.Queries instance for the given context.
// If the context contains a transaction, it returns queries bound to that transaction.
// Otherwise, it returns queries bound to the main database connection.
func (s *SQLiteStore) conn(ctx context.Context) *queries.Queries {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return s.queries.WithTx(tx)
	}
	return s.queries
}
