package store

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

// InitDB initializes the database connection and applies the schema
func InitDB(dbPath string) (*sql.DB, *Queries, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Apply schema
	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to apply schema: %w", err)
	}

	// Create queries instance
	queries := New(db)

	return db, queries, nil
}
