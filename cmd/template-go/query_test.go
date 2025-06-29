package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDatabaseInfoWithInvalidDSN(t *testing.T) {
	ctx := context.Background()
	_, err := getDatabaseInfo(ctx, "invalid-dsn")
	require.Error(t, err)
	// The pgx driver may parse some invalid DSNs and fail on connection instead
	require.Error(t, err, "Should return an error for invalid DSN")
}

func TestGetDatabaseInfoWithTimeout(t *testing.T) {
	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for the context to be cancelled
	time.Sleep(1 * time.Millisecond)

	dsn := "postgres://user:pass@127.0.0.1:5432/db?sslmode=disable"
	_, err := getDatabaseInfo(ctx, dsn)
	require.Error(t, err)
}

func TestGetDatabaseInfoWithValidConnection(t *testing.T) {
	// This test uses the real container setup
	if testDSN == "" {
		t.Skip("Skipping test, no testcontainer DSN available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbInfo, err := getDatabaseInfo(ctx, testDSN)
	require.NoError(t, err)

	// Verify all expected fields are present
	assert.Contains(t, dbInfo, "database")
	assert.Contains(t, dbInfo, "version")

	// Verify field types
	database, ok := dbInfo["database"].(string)
	assert.True(t, ok, "database should be a string")
	assert.NotEmpty(t, database)

	version, ok := dbInfo["version"].(string)
	assert.True(t, ok, "version should be a string")
	assert.NotEmpty(t, version)
}

// TestContextTimeout tests that our function respects context timeouts.
func TestContextTimeoutHandling(t *testing.T) {
	if testDSN == "" {
		t.Skip("Skipping test, no testcontainer DSN available")
	}

	// Test that function adds timeout when none exists
	ctxWithoutTimeout := context.Background()
	start := time.Now()

	_, err := getDatabaseInfo(ctxWithoutTimeout, testDSN)
	elapsed := time.Since(start)

	// Should complete successfully and within reasonable time (our 5s timeout)
	require.NoError(t, err)
	assert.Less(t, elapsed, 6*time.Second, "should complete within timeout")
}

// Test helper to validate SQL connection behavior.
func TestSQLConnectionBehavior(t *testing.T) {
	if testDSN == "" {
		t.Skip("Skipping test, no testcontainer DSN available")
	}

	// Test that sql.Open doesn't actually connect
	db, err := sql.Open("pgx", testDSN)
	require.NoError(t, err, "sql.Open should not fail with valid DSN")
	defer db.Close()

	// Test that Ping actually tests the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	require.NoError(t, err, "Ping should succeed with valid connection")
}
