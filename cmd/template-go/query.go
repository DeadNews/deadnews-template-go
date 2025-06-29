package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// getDatabaseInfo connects to the database and returns current database name and version.
func getDatabaseInfo(ctx context.Context, dsn string) (map[string]interface{}, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get database name and version
	var dbName, version string

	err = db.QueryRowContext(ctx, "SELECT current_database()").Scan(&dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to get current database name: %w", err)
	}

	err = db.QueryRowContext(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("failed to get database version: %w", err)
	}

	return map[string]interface{}{
		"database": dbName,
		"version":  version,
	}, nil
}
