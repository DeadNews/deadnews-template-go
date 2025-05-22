package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// healthCheck sends a request to the health endpoint of the running service.
func healthCheck(url string) error {
	// Create a context with timeout to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a request with proper validation and context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status code: %d", resp.StatusCode)
	}

	return nil
}
