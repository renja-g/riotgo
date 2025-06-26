package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/renja-g/riotgo/clients"
	"github.com/renja-g/riotgo/middleware"
)

func TestMiddlewareIntegration(t *testing.T) {
	// 1. Setup a mock server to receive the request.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 5. Assert that the headers from the middleware were added.
		if r.Header.Get("X-Test-Header-1") != "value1" {
			t.Error("Expected 'X-Test-Header-1' to be 'value1'")
		}
		if r.Header.Get("X-Test-Header-2") != "value2" {
			t.Error("Expected 'X-Test-Header-2' to be 'value2'")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 2. Create two simple middlewares that add headers to the request.
	middleware1 := func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("X-Test-Header-1", "value1")
			return next.RoundTrip(req)
		})
	}
	middleware2 := func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("X-Test-Header-2", "value2")
			return next.RoundTrip(req)
		})
	}

	// 3. Create a new BaseClient with the middlewares.
	client := clients.NewBaseClient(
		server.URL,
		clients.WithMiddleware(middleware1),
		clients.WithMiddleware(middleware2),
	)

	// 4. Make a request using the client.
	_, err := client.Invoke(context.Background(), http.MethodGet, "/", nil, nil, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
