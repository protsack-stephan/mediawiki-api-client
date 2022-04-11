// Package mediawiki API client for accessing actions API and REST base.
package mediawiki

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ErrEmptyResult no items in result list.
var ErrEmptyResult = errors.New("empty response result")

// ErrPageNotFound page not found in the api.
var ErrPageNotFound = errors.New("page not found")

// ErrUserNotFound user was not found.
var ErrUserNotFound = errors.New("user not found")

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
			pageDataURL,
			userURL,
		},
	}
}

// Client wikimedia api client.
type Client struct {
	url        string
	httpClient *http.Client
	headers    map[string]string
	options    *Options
}

// PageMeta get page meta data.
func (cl *Client) PageMeta(ctx context.Context, title string) (*PageMeta, error) {
	meta := new(PageMeta)
	data, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.PageMetaURL+url.QueryEscape(title), nil, cl.headers)

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

// PagesData get page data from Actions API.
func (cl *Client) PagesData(ctx context.Context, titles []string, options ...PageDataOptions) (map[string]PageData, error) {
	var rvProps []string
	rvLimit := 1
	pages := make(map[string]PageData)
	res := new(pageDataResponse)

	for _, opt := range options {
		rvLimit = opt.RevisionsLimit
		rvProps = opt.RevisionProps
	}

	body := url.Values{
		"action":        []string{"query"},
		"prop":          []string{"info|categories|revisions|templates|wbentityusage|pageprops|redirects|flagged"},
		"rvprop":        []string{"comment|oresscores|content|ids|timestamp|tags|user|userid|flags"},
		"rvslots":       []string{"main"},
		"inprop":        []string{"displaytitle|protection|url|watchers"},
		"ppprop":        []string{"wikibase_item"},
		"redirects":     []string{"1"},
		"titles":        []string{strings.Join(titles, "|")},
		"formatversion": []string{"2"},
		"format":        []string{"json"},
		"rvlimit":       []string{fmt.Sprintf("%d", rvLimit)},
	}

	if len(rvProps) > 0 {
		body["rvprop"][0] += fmt.Sprintf("|%s", strings.Join(rvProps, "|"))
	}

	data, status, err := req(
		ctx,
		cl.httpClient,
		http.MethodPost,
		fmt.Sprintf("%s%s", cl.url, cl.options.PageDataURL),
		strings.NewReader(body.Encode()),
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}, cl.headers)

	if err != nil {
		return pages, err
	}

	if status != http.StatusOK {
		return pages, fmt.Errorf(errBadRequestMsg, status, data)
	}

	if err := json.Unmarshal(data, res); err != nil {
		return pages, err
	}

	lookup := map[string]bool{}

	for _, title := range titles {
		lookup[title] = true
	}

	normalized := map[string]string{}

	for _, title := range res.Query.Normalized {
		if _, ok := lookup[title.From]; ok {
			normalized[title.To] = title.From
		}
	}

	for _, page := range res.Query.Pages {
		if !page.Missing {
			if title, ok := normalized[page.Title]; ok {
				pages[title] = page
			} else if _, ok := lookup[page.Title]; ok {
				pages[page.Title] = page
			}
		}
	}

	return pages, nil
}

// PageData get page data from Actions API.
func (cl *Client) PageData(ctx context.Context, title string, options ...PageDataOptions) (PageData, error) {
	resp, err := cl.PagesData(ctx, []string{title}, options...)

	if err != nil {
		return resp[title], err
	}

	if data, ok := resp[title]; ok {
		return data, nil
	}

	return resp[title], ErrPageNotFound
}

// PageHTML get page html with or without revision.
func (cl *Client) PageHTML(ctx context.Context, title string, rev ...int) ([]byte, error) {
	url := cl.url + cl.options.PageHTMLURL + url.QueryEscape(title)

	if len(rev) > 0 {
		url += "/" + strconv.Itoa(rev[0])
	}

	data, status, err := req(ctx, cl.httpClient, http.MethodGet, url, nil, cl.headers)

	if err != nil {
		return []byte{}, err
	}

	if status != http.StatusOK {
		return []byte{}, fmt.Errorf(errBadRequestMsg, status, data)
	}

	return data, nil
}

// PageWikitext get page wikitext with or without revision.
func (cl *Client) PageWikitext(ctx context.Context, title string, rev ...int) ([]byte, error) {
	url := cl.url + fmt.Sprintf(cl.options.PageWikitextURL, url.QueryEscape(title))

	if len(rev) > 0 {
		url += "&rvstartid=" + strconv.Itoa(rev[0])
	}

	data, status, err := req(ctx, cl.httpClient, http.MethodGet, url, nil, cl.headers)

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

// PageRevisions get list of page revisions.
func (cl *Client) PageRevisions(ctx context.Context, title string, limit int, options ...PageRevisionsOptions) ([]Revision, error) {
	var props []string
	revs := []Revision{}
	ordering := RevisionOrderingOlder

	for _, opt := range options {
		ordering = opt.Order
		props = opt.Props
	}

	body := url.Values{
		"action":        []string{"query"},
		"format":        []string{"json"},
		"formatversion": []string{"2"},
		"prop":          []string{"revisions"},
		"titles":        []string{url.QueryEscape(title)},
		"rvlimit":       []string{strconv.Itoa(limit)},
		"rvdir":         []string{string(ordering)},
	}

	if len(props) > 0 {
		body["rvprop"] = []string{strings.Join(props, "|")}
	}

	reqURL := cl.url + cl.options.PageRevisionsURL + body.Encode()
	data, status, err := req(ctx, cl.httpClient, http.MethodGet, reqURL, nil, cl.headers)

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

// Sitematrix get all supported wikimedia projects.
func (cl *Client) Sitematrix(ctx context.Context) (*Sitematrix, error) {
	matrix := new(Sitematrix)
	data, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.SitematrixURL, nil, cl.headers)

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

// Namespaces get page types called "namespaces".
func (cl *Client) Namespaces(ctx context.Context) ([]Namespace, error) {
	ns := []Namespace{}
	data, status, err := req(ctx, cl.httpClient, http.MethodGet, cl.url+cl.options.NamespacesURL, nil, cl.headers)

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

// Users get list of users by id.
func (cl *Client) Users(ctx context.Context, ids ...int) (map[int]User, error) {
	ususerids := []string{}

	for _, id := range ids {
		ususerids = append(ususerids, strconv.Itoa(id))
	}

	body := url.Values{
		"action":        []string{"query"},
		"list":          []string{"users"},
		"usprop":        []string{"groups|editcount|groupmemberships|registration|emailable"},
		"format":        []string{"json"},
		"formatversion": []string{"2"},
		"ususerids":     []string{strings.Join(ususerids, "|")},
	}

	data, status, err := req(
		ctx,
		cl.httpClient,
		http.MethodPost,
		fmt.Sprintf("%s%s", cl.url, cl.options.UserURL),
		strings.NewReader(body.Encode()),
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}, cl.headers)

	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf(errBadRequestMsg, status, data)
	}

	res := new(userResponse)

	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}

	users := make(map[int]User)

	for _, user := range res.Query.Users {
		if !user.Missing {
			users[user.UserID] = user
		}
	}

	return users, nil
}

// Users get single user by id.
func (cl *Client) User(ctx context.Context, id int) (User, error) {
	users, err := cl.Users(ctx, id)

	if err != nil {
		return users[id], err
	}

	if data, ok := users[id]; ok {
		return data, nil
	}

	return users[id], ErrUserNotFound
}
