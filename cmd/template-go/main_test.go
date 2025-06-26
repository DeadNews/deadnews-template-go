package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeServer(t *testing.T) {
	server := makeServer(":8080")
	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.Addr)
	assert.NotNil(t, server.Handler)
}

func TestServerEndpoints(t *testing.T) {
	server := makeServer(":0")
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		path     string
		expected string
		content  string
	}{
		{"/test", `{"status":"ok"}` + "\n", "application/json"},
		{"/health", ".", "text/plain"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+tt.path, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Logf("Failed to close response body: %v", err)
				}
			}()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Contains(t, resp.Header.Get("Content-Type"), tt.content)
		})
	}
}
