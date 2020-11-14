package mediawiki

import (
	"net/http"
	"time"
)

// NewClient create new wikimedia api client
func NewClient(url string) *Client {
	return &Client{
		url:        url,
		httpClient: new(http.Client),
		options: &Options{
			PageMetaURL: pageMetaURL,
			PageHTMLURL: pageHTMLURL,
		},
	}
}

type clientBuilder struct {
	client *Client
}

// Options update the default options
func (cb *clientBuilder) Options(options *Options) *clientBuilder {
	cb.client.options = options
	return cb
}

// HTTPClient pass custom http client
func (cb *clientBuilder) HTTPClient(httpClient *http.Client) *clientBuilder {
	cb.client.httpClient = httpClient
	return cb
}

func (cb *clientBuilder) Timeout(timeout time.Duration) *clientBuilder {
	cb.client.httpClient.Timeout = timeout
	return cb
}
