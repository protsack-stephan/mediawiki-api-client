package mediawiki

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ErrEmptyResult not items in result list
var ErrEmptyResult = errors.New("empty response result")

const errBadRequestMsg = "status: '%d' body: '%s'"

// NewClient create new client instance
func NewClient(url string) *Client {
	return &Client{
		url:        url,
		httpClient: new(http.Client),
		options: &Options{
			pageMetaURL,
			pageHTMLURL,
			pageWikitextURL,
			revisionsURL,
			sitematrixURL,
			namespacesURL,
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
func (cl *Client) PageMeta(ctx context.Context, title string) (*PageMeta, error) {
	meta := new(PageMeta)
	data, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.PageMetaURL+url.QueryEscape(title), nil)

	if err != nil {
		return meta, err
	}

	if status != http.StatusOK {
		return meta, fmt.Errorf(errBadRequestMsg, status, data)
	}

	res := new(pageMetaResponse)
	err = json.Unmarshal(data, res)

	if err != nil {
		return meta, err
	}

	if res.Items == nil || len(res.Items) <= 0 {
		return meta, ErrEmptyResult
	}

	return &res.Items[0], nil
}

// PageHTML get page html with with or without revision
func (cl *Client) PageHTML(ctx context.Context, title string, rev ...int) ([]byte, error) {
	url := cl.url + cl.options.PageHTMLURL + url.QueryEscape(title)

	if len(rev) > 0 {
		url += "/" + strconv.Itoa(rev[0])
	}

	data, status, err := req(ctx, cl.httpClient, http.MethodGet, url, nil)

	if err != nil {
		return []byte{}, err
	}

	if status != http.StatusOK {
		return []byte{}, fmt.Errorf(errBadRequestMsg, status, data)
	}

	return data, nil
}

// PageWikitext get page wikitext with or without revision
func (cl *Client) PageWikitext(ctx context.Context, title string, rev ...int) ([]byte, error) {
	url := cl.url + fmt.Sprintf(cl.options.PageWikitextURL, url.QueryEscape(title))

	if len(rev) > 0 {
		url += "&rvstartid=" + strconv.Itoa(rev[0])
	}

	data, status, err := req(ctx, cl.httpClient, http.MethodGet, url, nil)

	if err != nil {
		return []byte{}, err
	}

	if status != http.StatusOK {
		return []byte{}, fmt.Errorf(errBadRequestMsg, status, data)
	}

	res := new(wikitextResponse)
	err = json.Unmarshal(data, res)

	if err != nil {
		return []byte{}, err
	}

	if len(res.Query.Pages) <= 0 || len(res.Query.Pages[0].Revisions) <= 0 {
		return []byte{}, ErrEmptyResult
	}

	return []byte(res.Query.Pages[0].Revisions[0].Slots.Main.Content), err
}

// PageRevisions get list of page revisions
func (cl *Client) PageRevisions(ctx context.Context, title string, limit int) ([]Revision, error) {
	revs := []Revision{}
	data, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+fmt.Sprintf(cl.options.PageRevisionsURL, limit, url.QueryEscape(title)), nil)

	if err != nil {
		return revs, err
	}

	if status != http.StatusOK {
		return revs, fmt.Errorf(errBadRequestMsg, status, data)
	}

	res := new(revisionsResponse)
	err = json.Unmarshal(data, res)

	if err != nil {
		return revs, err
	}

	if len(res.Query.Pages) == 0 || len(res.Query.Pages[0].Revisions) == 0 {
		return revs, ErrEmptyResult
	}

	return res.Query.Pages[0].Revisions, nil
}

// Sitematrix get all supported wikimedia projects
func (cl *Client) Sitematrix(ctx context.Context) (*Sitematrix, error) {
	matrix := new(Sitematrix)
	data, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.SitematrixURL, nil)

	if err != nil {
		return matrix, err
	}

	if status != http.StatusOK {
		return matrix, fmt.Errorf(errBadRequestMsg, status, data)
	}

	special := new(siteMatrixSpecialResponce)
	err = json.Unmarshal(data, special)

	if err != nil {
		return matrix, err
	}

	matrix.Count, matrix.Specials = special.Sitematrix.Count, special.Sitematrix.Specials
	main := new(siteMatrixMainResponse)
	_ = json.Unmarshal(data, main)

	for num, project := range main.Sitematrix {
		if num != "count" && num != "specials" {
			matrix.Projects = append(matrix.Projects, project)
		}
	}

	return matrix, nil
}

// Namespaces get page types called "namespaces"
func (cl *Client) Namespaces(ctx context.Context) ([]Namespace, error) {
	ns := []Namespace{}
	data, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.NamespacesURL, nil)

	if err != nil {
		return ns, err
	}

	if status != http.StatusOK {
		return ns, fmt.Errorf(errBadRequestMsg, status, data)
	}

	res := new(namespacesResponse)
	err = json.Unmarshal(data, res)

	if err != nil {
		return ns, err
	}

	for _, name := range res.Query.Namespaces {
		ns = append(ns, name)
	}

	return ns, nil
}
