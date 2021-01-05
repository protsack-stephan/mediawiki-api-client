package mediawiki

import "time"

const pageDataURL = "/w/api.php?action=query&prop=pageprops|info|revisions&ppprop=wikibase_item&redirects=1&titles=%s&formatversion=2&format=json"

// PageData page data returned from actions API
type PageData struct {
	PageID    int    `json:"pageid"`
	Ns        int    `json:"ns"`
	Title     string `json:"title"`
	Missing   bool   `json:"missing"`
	Pageprops struct {
		WikibaseItem string `json:"wikibase_item"`
	} `json:"pageprops"`
	ContentModel         string    `json:"contentmodel"`
	PageLanguage         string    `json:"pagelanguage"`
	PageLanguageHTMLCode string    `json:"pagelanguagehtmlcode"`
	PageLanguageDir      string    `json:"pagelanguagedir"`
	Touched              time.Time `json:"touched"`
	LastRevID            int       `json:"lastrevid"`
	Length               int       `json:"length"`
	Revisions            []struct {
		RevID     int       `json:"revid"`
		ParentID  int       `json:"parentid"`
		Minor     bool      `json:"minor"`
		User      string    `json:"user"`
		Timestamp time.Time `json:"timestamp"`
		Comment   string    `json:"comment"`
	} `json:"revisions"`
}

type pageDataResponse struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Pages []PageData `json:"pages"`
	} `json:"query"`
}
