package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLFor(t *testing.T) {
	t.Parallel()
	rc := &RiotAPIClient{BaseClient: NewBaseClient("https://%s.api.riotgames.com")}
	got := rc.urlFor(Europe, "/v1/test")
	want := "https://europe.api.riotgames.com/v1/test"
	if got != want {
		t.Fatalf("urlFor() = %s, want %s", got, want)
	}
}

func TestInvokeJSONSuccess(t *testing.T) {
	t.Parallel()

	type resp struct {
		Msg string `json:"msg"`
	}

	// stub server returns 200 JSON
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp{Msg: "ok"})
	}))
	defer srv.Close()

	// Create client pointing to stub (no %s spec) so region arg is ignored
	rc := &RiotAPIClient{BaseClient: NewBaseClient("%s", WithHTTPClient(srv.Client()))}

	out, err := invokeJSON[resp](rc, context.Background(), Region(srv.URL), http.MethodGet, "/foo", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("invokeJSON returned error: %v", err)
	}
	if out.Msg != "ok" {
		t.Errorf("unexpected payload: %+v", out)
	}
}

func TestInvokeJSONErrorStatus(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusTeapot)
	}))
	defer srv.Close()

	rc := &RiotAPIClient{BaseClient: NewBaseClient("%s", WithHTTPClient(srv.Client()))}

	_, err := invokeJSON[struct{}](rc, context.Background(), Region(srv.URL), http.MethodGet, "/foo", nil, nil, nil, nil)
	if err == nil {
		t.Fatalf("expected error for non-2xx status")
	}
}
