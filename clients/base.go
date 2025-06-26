package clients

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/renja-g/riotgo/middleware"
)

type BaseClient struct {
	BaseURL            string
	DefaultHeaders     map[string]string
	DefaultParams      map[string]string
	DefaultQueryParams map[string]string
	HTTPClient         *http.Client
	middleware         []middleware.Middleware
}

func NewBaseClient(baseURL string, opts ...Option) *BaseClient {
	c := &BaseClient{
		BaseURL:            strings.TrimRight(baseURL, "/"),
		DefaultHeaders:     map[string]string{},
		DefaultParams:      map[string]string{},
		DefaultQueryParams: map[string]string{},
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Option func(*BaseClient)

func WithHTTPClient(h *http.Client) Option {
	return func(c *BaseClient) { c.HTTPClient = h }
}

func WithDefaultHeader(k, v string) Option {
	return func(c *BaseClient) { c.DefaultHeaders[k] = v }
}

func WithDefaultQuery(k, v string) Option {
	return func(c *BaseClient) { c.DefaultQueryParams[k] = v }
}

func WithMiddleware(m middleware.Middleware) Option {
	return func(c *BaseClient) { c.middleware = append(c.middleware, m) }
}

func (c *BaseClient) Invoke(
	ctx context.Context,
	method, path string,
	body io.Reader,
	headers map[string]string,
	queries url.Values,
) (*http.Response, error) {

	var u string
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		u = path
	} else {
		u = c.BaseURL + path
	}

	q := url.Values{}
	for k, v := range c.DefaultQueryParams {
		q.Set(k, v)
	}
	for k, vs := range queries {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	if encoded := q.Encode(); encoded != "" {
		u += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.DefaultHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := *c.HTTPClient
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	// Chain the middleware in reverse order.
	for i := len(c.middleware) - 1; i >= 0; i-- {
		transport = c.middleware[i](transport)
	}
	client.Transport = transport

	return client.Do(req)
}
