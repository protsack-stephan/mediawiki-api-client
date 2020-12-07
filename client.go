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
func (cl *Client) PageHTML(ctx context.Context, title string, rev ...int) ([]byte, int, error) {
	url := cl.url + cl.options.PageHTMLURL + url.QueryEscape(title)

	if len(rev) > 0 {
		url += "/" + strconv.Itoa(rev[0])
	}

	return req(ctx, cl.httpClient, http.MethodGet, url, nil)
}

// PageWikitext get page wikitext with or without revision
func (cl *Client) PageWikitext(ctx context.Context, title string, rev ...int) ([]byte, int, error) {
	url := cl.url + fmt.Sprintf(cl.options.PageWikitextURL, url.QueryEscape(title))

	if len(rev) > 0 {
		url += "&rvstartid=" + strconv.Itoa(rev[0])
	}

	res, status, err := req(ctx, cl.httpClient, http.MethodGet, url, nil)

	if err != nil {
		return []byte{}, status, err
	}

	wRes := new(wikitextResponse)
	err = json.Unmarshal(res, wRes)

	if err != nil {
		return []byte{}, status, err
	}

	if len(wRes.Query.Pages) <= 0 || len(wRes.Query.Pages[0].Revisions) <= 0 {
		return []byte{}, status, fmt.Errorf("no data in result")
	}

	return []byte(wRes.Query.Pages[0].Revisions[0].Slots.Main.Content), status, err
}

// PageRevisions get list of page revisions
func (cl *Client) PageRevisions(ctx context.Context, title string, limit int) ([]Revision, int, error) {
	revs := []Revision{}
	res, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+fmt.Sprintf(cl.options.PageRevisionsURL, limit, url.QueryEscape(title)), nil)

	if err != nil {
		return revs, status, err
	}

	rRes := new(revisionsResponse)
	err = json.Unmarshal(res, rRes)

	if err != nil {
		return revs, status, err
	}

	if len(rRes.Query.Pages) == 0 {
		return revs, status, fmt.Errorf("revisions not found")
	}

	return rRes.Query.Pages[0].Revisions, status, nil
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
