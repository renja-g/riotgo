package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestExpandPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		template string
		params   map[string]string
		expect   string
	}{
		{
			name:     "single param",
			template: "/foo/{id}",
			params:   map[string]string{"id": "123"},
			expect:   "/foo/123",
		},
		{
			name:     "multiple params",
			template: "/foo/{id}/bar/{name}",
			params:   map[string]string{"id": "1", "name": "baz"},
			expect:   "/foo/1/bar/baz",
		},
		{
			name:     "missing param untouched",
			template: "/foo/{id}/bar/{name}",
			params:   map[string]string{"id": "9"},
			expect:   "/foo/9/bar/{name}", // placeholders without value remain
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expandPath(tt.template, tt.params); got != tt.expect {
				t.Errorf("expandPath() = %s, want %s", got, tt.expect)
			}
		})
	}
}

func TestNewBaseClientOptions(t *testing.T) {
	t.Parallel()

	customHTTP := &http.Client{Timeout: 42 * time.Second}
	bc := NewBaseClient("http://example.com",
		WithHTTPClient(customHTTP),
		WithDefaultHeader("X-Foo", "Bar"),
		WithDefaultQuery("a", "b"),
	)

	if bc.HTTPClient != customHTTP {
		t.Errorf("expected custom HTTP client to be set")
	}
	if v := bc.DefaultHeaders["X-Foo"]; v != "Bar" {
		t.Errorf("default header not set, got %q", v)
	}
	if v := bc.DefaultQueryParams["a"]; v != "b" {
		t.Errorf("default query param not set, got %q", v)
	}
}

func TestInvokeURLBuilding(t *testing.T) {
	t.Parallel()

	var capturedURL *url.URL
	var capturedHeader http.Header

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL
		capturedHeader = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	bc := NewBaseClient(srv.URL,
		WithHTTPClient(srv.Client()),
		WithDefaultHeader("X-Default", "D"),
		WithDefaultQuery("def", "1"),
	)

	queries := url.Values{"extra": []string{"2"}}
	headers := map[string]string{"X-Custom": "C"}

	if _, err := bc.Invoke(context.Background(), http.MethodGet, "/foo", nil, headers, queries); err != nil {
		t.Fatalf("Invoke returned error: %v", err)
	}

	expectedPath := "/foo?def=1&extra=2"
	if capturedURL == nil {
		t.Fatalf("no request captured")
	}
	if capturedURL.String() != expectedPath {
		t.Errorf("url mismatch: got %s want %s", capturedURL.String(), expectedPath)
	}

	if hv := capturedHeader.Get("X-Default"); hv != "D" {
		t.Errorf("default header missing, got %q", hv)
	}
	if hv := capturedHeader.Get("X-Custom"); hv != "C" {
		t.Errorf("custom header missing, got %q", hv)
	}
}
