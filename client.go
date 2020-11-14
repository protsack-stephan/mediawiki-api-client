package mediawiki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// NewClient create new wikimedia api client
func NewClient(url string) *Client {
	return &Client{
		url:        url,
		httpClient: new(http.Client),
		options: &Options{
			PageMetaURL: pageMetaURL,
		},
	}
}

// Client wikimedia api client
type Client struct {
	url        string
	httpClient *http.Client
	options    *Options
}

// Options update the default options
func (cl *Client) Options(options *Options) *Client {
	cl.options = options
	return cl
}

// HTTPClient pass custom http client
func (cl *Client) HTTPClient(httpClient *http.Client) *Client {
	cl.httpClient = httpClient
	return cl
}

// PageMeta get page meta data
func (cl *Client) PageMeta(ctx context.Context, title string) (*PageMeta, int, error) {
	meta := new(PageMetaResponse)
	res, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.PageMetaURL+url.QueryEscape(title), nil)

	if err != nil {
		return nil, status, err
	}

	err = json.Unmarshal(res, &meta)

	if err != nil {
		return nil, status, err
	}

	if meta.Items == nil || len(meta.Items) <= 0 {
		return nil, status, fmt.Errorf("zero items in result")
	}

	return &meta.Items[0], status, nil
}
