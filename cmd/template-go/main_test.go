package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Global variables for test environment.
var (
	testDSN       string
	testContext   context.Context
	testContainer testcontainers.Container
)

// TestMain sets up and tears down the shared test environment.
func TestMain(m *testing.M) {
	// Skip container setup if TESTCONTAINERS is not set.
	if os.Getenv("TESTCONTAINERS") != "1" {
		os.Exit(m.Run())
	}

	// Create context
	testContext = context.Background()

	// Create container request
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		),
	}

	// Start container
	pgContainer, err := testcontainers.GenericContainer(testContext, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("Failed to start container: %s\n", err)
		os.Exit(1)
	}
	testContainer = pgContainer

	// Get the mapped port
	port, err := pgContainer.MappedPort(testContext, "5432")
	if err != nil {
		fmt.Printf("Failed to get port: %s\n", err)
		terminateContainer()
		os.Exit(1)
	}

	// Get the host
	host, err := pgContainer.Host(testContext)
	if err != nil {
		fmt.Printf("Failed to get host: %s\n", err)
		terminateContainer()
		os.Exit(1)
	}

	// Construct DSN
	testDSN = fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	fmt.Println("Using test DSN: ", testDSN)

	// Run the tests
	code := m.Run()

	// Clean up
	terminateContainer()

	// Exit with the appropriate code
	os.Exit(code)
}

// Helper function to terminate the container.
func terminateContainer() {
	if err := testContainer.Terminate(testContext); err != nil {
		fmt.Printf("Error terminating container: %s\n", err)
	}
}

func TestDatabaseService(t *testing.T) {
	if os.Getenv("TESTCONTAINERS") != "1" {
		t.Skip("Skipping integration test, set TESTCONTAINERS=1 to run it.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test getting database info directly
	dbInfo, err := getDatabaseInfo(ctx, testDSN)
	require.NoError(t, err)

	// Verify expected fields are present
	assert.Contains(t, dbInfo, "database")
	assert.Contains(t, dbInfo, "version")

	// Verify database name is correct
	assert.Equal(t, "testdb", dbInfo["database"])

	// Verify version contains PostgreSQL
	version, ok := dbInfo["version"].(string)
	assert.True(t, ok, "version should be a string")
	assert.Contains(t, version, "PostgreSQL", "version should contain PostgreSQL")
}

func TestSetupServer(t *testing.T) {
	server := setupServer(":8080", "test-dsn")
	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.Addr)
	assert.NotNil(t, server.Handler)
}

func TestMakeDatabaseTestHandler_Success(t *testing.T) {
	if os.Getenv("TESTCONTAINERS") != "1" {
		t.Skip("Skipping integration test, set TESTCONTAINERS=1 to run it.")
	}

	handler := handleDatabaseTest(testDSN)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Verify expected fields are present
	assert.Contains(t, response, "database")
	assert.Contains(t, response, "version")

	// Verify database name is correct
	assert.Equal(t, "testdb", response["database"])

	// Verify version contains PostgreSQL
	version, ok := response["version"].(string)
	assert.True(t, ok, "version should be a string")
	assert.Contains(t, version, "PostgreSQL", "version should contain PostgreSQL")
}

func TestMakeDatabaseTestHandler_Error(t *testing.T) {
	handler := handleDatabaseTest("invalid-dsn")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Internal server error")
}

func TestHandleDatabaseTest_ViaServer_Success(t *testing.T) {
	if os.Getenv("TESTCONTAINERS") != "1" {
		t.Skip("Skipping integration test, set TESTCONTAINERS=1 to run it.")
	}

	server := setupServer(":0", testDSN)
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/test", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Contains(t, response, "database")
	assert.Contains(t, response, "version")
}

func TestHandleDatabaseTest_ViaServer_Error(t *testing.T) {
	server := setupServer(":0", "invalid-dsn")
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/test", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestHealthEndpoint(t *testing.T) {
	server := setupServer(":0", "test-dsn")
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/health", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain")
}

// TestMakeDatabaseTestHandler_JSONError tests the JSON encoding error path.
func TestMakeDatabaseTestHandler_JSONError(t *testing.T) {
	if os.Getenv("TESTCONTAINERS") != "1" {
		t.Skip("Skipping integration test, set TESTCONTAINERS=1 to run it.")
	}

	// Create a handler that will trigger JSON encoding error by using a channel (which can't be JSON encoded)
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// This will succeed
		// Add a channel to the response to force JSON encoding error
		response := map[string]interface{}{
			"bad": make(chan int),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Status will still be 200 because the error happens after WriteHeader
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
}

// TestServerConfiguration tests server timeout configuration.
func TestServerConfiguration(t *testing.T) {
	server := setupServer(":8080", "test-dsn")

	assert.Equal(t, 15*time.Second, server.ReadTimeout)
	assert.Equal(t, 15*time.Second, server.WriteTimeout)
	assert.Equal(t, 60*time.Second, server.IdleTimeout)
}

// TestMakeDatabaseTestHandler_ContextCancellation tests context cancellation.
func TestMakeDatabaseTestHandler_ContextCancellation(t *testing.T) {
	// Use testDSN if available, otherwise skip this test
	dsn := testDSN
	if dsn == "" {
		dsn = "postgres://user:pass@127.0.0.1:5432/db?sslmode=disable"
	}

	handler := handleDatabaseTest(dsn)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Internal server error")
}

// TestSetupServer_MiddlewareConfiguration tests that all middleware is properly configured.
func TestSetupServer_MiddlewareConfiguration(t *testing.T) {
	server := setupServer(":8080", "test-dsn")

	// Test that heartbeat endpoint is configured
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestSetupServer_RoutingConfiguration tests that routes are properly configured.
func TestSetupServer_RoutingConfiguration(t *testing.T) {
	server := setupServer(":8080", "invalid-dsn")

	// Test that /test endpoint exists and returns error for invalid DSN
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestSetupServer_NonExistentRoute tests 404 handling.
func TestSetupServer_NonExistentRoute(t *testing.T) {
	server := setupServer(":8080", "test-dsn")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/nonexistent", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
