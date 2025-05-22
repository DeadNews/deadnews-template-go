// Package main is the entry point for the template-go application.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// This code is the main function of a Go program
// that creates and starts a server using the Echo framework.
// It retrieves the value of the "SERVICE_PORT" environment variable
// and if it is not set, it defaults to port 8000.
func main() {
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

	// Create a new Echo instance.
	e := makeServer()

	// Get the value of the "SERVICE_PORT" environment variable.
	httpPort := os.Getenv("SERVICE_PORT")
	if httpPort == "" {
		httpPort = "8000"
	}

	// Start the server on the specified port.
	e.Logger.Fatal(e.Start(":" + httpPort))
}

// makeServer creates a new instance of the Echo framework
// and configures it with middleware for logging and error recovery.
// It also defines two route handlers: one for the root ("/") route that returns an HTML response,
// and another for the "/health" route that returns a JSON response.
//
// Returns:
// - A configured instance of the Echo framework.
func makeServer() *echo.Echo {
	// Create a new Echo instance.
	e := echo.New()

	// Use middleware for logging and error recovery.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the "/" route handler.
	e.GET("/", handleRoot)

	// Define the "/health" route handler.
	e.GET("/health", handleHealth)

	return e
}

// handleRoot handles the "/" route and returns an HTML response.
func handleRoot(c echo.Context) error {
	return c.HTML(http.StatusOK, "Hello, World!\n")
}

// handleHealth handles the "/health" route and returns a JSON response.
func handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
}
