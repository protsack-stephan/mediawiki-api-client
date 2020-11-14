package mediawiki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Client wikimedia api client
type Client struct {
	url        string
	httpClient *http.Client
	options    *Options
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

// PageHTML get page html with revision
func (cl *Client) PageHTML(ctx context.Context, title string, rev ...int) ([]byte, int, error) {
	url := cl.url + cl.options.PageHTMLURL + url.QueryEscape(title)

	if len(rev) > 0 {
		url += "/" + strconv.Itoa(rev[0])
	}

	return req(ctx, cl.httpClient, http.MethodGet, url, nil)
}
