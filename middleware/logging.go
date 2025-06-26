package middleware

import (
	"log"
	"net/http"
)

func Logging() Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			log.Printf("--> %s %s", req.Method, req.URL)
			resp, err := next.RoundTrip(req)
			if err == nil {
				log.Printf("<-- %s", resp.Status)
			}
			return resp, err
		})
	}
}
