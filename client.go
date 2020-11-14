package mediawiki

import "net/http"

// NewClient create new wikimedia api client
func NewClient(url string) *Client {
	return &Client{
		URL:  url,
		HTTP: new(http.Client),
	}
}

// Client wikimedia api client
type Client struct {
	URL  string
	HTTP *http.Client
}


func (cl* Client) PageMeta()