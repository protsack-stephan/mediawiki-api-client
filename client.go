package mediawiki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// NewClient create new client instance
func NewClient(url string) *Client {
	return &Client{
		url:        url,
		httpClient: new(http.Client),
		options: &Options{
			PageMetaURL:   pageMetaURL,
			PageHTMLURL:   pageHTMLURL,
			SitematrixURL: sitematrixURL,
			NamespacesURL: namespacesURL,
		},
	}
}

// Client wikimedia api client
type Client struct {
	url        string
	httpClient *http.Client
	options    *Options
}

// PageMeta get page meta data
func (cl *Client) PageMeta(ctx context.Context, title string) (*PageMeta, int, error) {
	meta := new(pageMetaResponse)
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

// Sitematrix get all supported wikimedia projects
func (cl *Client) Sitematrix(ctx context.Context) (*Sitematrix, int, error) {
	matrix := new(Sitematrix)
	res, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.SitematrixURL, nil)

	if err != nil {
		return matrix, status, err
	}

	spRes := new(siteMatrixSpecialResponce)
	err = json.Unmarshal(res, spRes)

	if err != nil {
		return matrix, status, err
	}

	matrix.Count, matrix.Specials = spRes.Sitematrix.Count, spRes.Sitematrix.Specials
	mRes := new(siteMatrixMainResponse)
	json.Unmarshal(res, mRes)

	for num, project := range mRes.Sitematrix {
		if num != "count" && num != "specials" {
			matrix.Projects = append(matrix.Projects, project)
		}
	}

	return matrix, status, nil
}

// Namespaces get page types called "namespaces"
func (cl *Client) Namespaces(ctx context.Context) ([]Namespace, int, error) {
	ns := []Namespace{}
	res, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.NamespacesURL, nil)

	if err != nil {
		return ns, status, err
	}

	nsRes := new(namespacesResponse)
	err = json.Unmarshal(res, &nsRes)

	if err != nil {
		return ns, status, err
	}

	for _, name := range nsRes.Query.Namespaces {
		ns = append(ns, name)
	}

	return ns, status, nil
}
