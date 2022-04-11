package mediawiki

import (
	"net/http"
	"time"
)

// NewBuilder create new wikimedia api client
func NewBuilder(url string) *ClientBuilder {
	cb := new(ClientBuilder)
	cb.client = NewClient(url)
	return cb
}

// ClientBuilder builder method fo the client
type ClientBuilder struct {
	client *Client
}

// Options update the default options
func (cb *ClientBuilder) Options(options *Options) *ClientBuilder {
	cb.client.options = options
	return cb
}

// HTTPClient pass custom http client
func (cb *ClientBuilder) HTTPClient(httpClient *http.Client) *ClientBuilder {
	cb.client.httpClient = httpClient
	return cb
}

// Timeout set http client timeout
func (cb *ClientBuilder) Timeout(timeout time.Duration) *ClientBuilder {
	cb.client.httpClient.Timeout = timeout
	return cb
}

// Headers set http headers.
func (cb *ClientBuilder) Headers(headers map[string]string) *ClientBuilder {
	cb.client.headers = headers
	return cb
}

// Build create new client instance
func (cb *ClientBuilder) Build() *Client {
	return cb.client
}
