// Package main is the entry point for the template-go application.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

func main() {
	// Setup structured logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Parse command-line flags
	healthcheckURL := flag.String("healthcheck", "", "Perform a health check against the given URL and exit")
	flag.Parse()

	// Handle health check mode
	if *healthcheckURL != "" {
		if err := healthCheck(*healthcheckURL); err != nil {
			fmt.Fprintf(os.Stderr, "Health check failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Health check succeeded")
		os.Exit(0)
	}

	// Get port from environment
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8000"
	}

	// Create and start server
	server := makeServer(":" + port)
	slog.Info("Starting server", "port", port)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

// makeServer creates a configured HTTP server with Chi router.
func makeServer(addr string) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	r.Get("/test", handleTest)

	return &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// handleTest returns a JSON health status.
func handleTest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}
