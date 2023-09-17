package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerCreation(t *testing.T) {
	e := makeServer()

	assert.NotNil(t, e.GET("/", handleRoot))
	assert.NotNil(t, e.GET("/health", handleHealth))
}

func TestServerResponseRoot(t *testing.T) {
	// Create a new Echo instance.
	e := makeServer()

	// Create a new HTTP request to the root route.
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create a new HTTP response recorder.
	rec := httptest.NewRecorder()

	// Handle the request.
	e.ServeHTTP(rec, req)

	// Check the actual status code against the expected status code.
	expectedStatus := http.StatusOK
	assert.Equal(t, expectedStatus, rec.Code)

	// Check the actual response body against the expected response body.
	expectedBody := "Hello, World!"
	assert.Equal(t, expectedBody, rec.Body.String())

	// Check the response content type header.
	contentType := "text/html; charset=UTF-8"
	assert.Equal(t, contentType, rec.Header().Get("Content-Type"))
}

func TestServerResponseHealth(t *testing.T) {
	// Create a new Echo instance.
	e := makeServer()

	// Create a new HTTP request with the "/health" route.
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Create a new HTTP response recorder.
	rec := httptest.NewRecorder()

	// Handle the request.
	e.ServeHTTP(rec, req)

	// Check the actual status code against the expected status code.
	expectedStatus := http.StatusOK
	assert.Equal(t, expectedStatus, rec.Code)

	// Check the actual response body against the expected response body.
	expectedBody := `{"Status":"OK"}` + "\n"
	assert.Equal(t, expectedBody, rec.Body.String())

	// Check the response content type header.
	contentType := "application/json; charset=UTF-8"
	assert.Equal(t, contentType, rec.Header().Get("Content-Type"))
}
