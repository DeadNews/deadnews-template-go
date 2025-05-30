package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		errSubstr string
	}{
		{
			name: "successful health check",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr: false,
		},
		{
			name: "service unavailable",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
			},
			wantErr:   true,
			errSubstr: "status code: 503",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server with the handler
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			// Run the health check against the test server
			err := healthCheck(server.URL)

			// Check results
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("error %q doesn't contain %q", err.Error(), tt.errSubstr)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
