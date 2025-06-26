package middleware

import (
	"net/http"
)

type Middleware = func(http.RoundTripper) http.RoundTripper

type RoundTripperFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
