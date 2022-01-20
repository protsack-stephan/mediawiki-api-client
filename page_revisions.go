package mediawiki

import "time"

const revisionsURL = "/w/api.php?action=query&format=json&formatversion=2&prop=revisions&rvlimit=%d&rvdir=%s&titles=%s"

// PageRevisionsOptions additional optional parameters for PageRevisions method
type PageRevisionsOptions struct {
	Order RevisionOrdering
	Props []string
}

// Revision page revision schema
type Revision struct {
	RevID     int       `json:"revid"`
	ParentID  int       `json:"parentid"`
	Minor     bool      `json:"minor"`
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
	Comment   string    `json:"comment"`
	Anon      bool      `json:"anon,omitempty"`
}

// RevisionOrdering specifies the direction to enumerate the revisions list
type RevisionOrdering string

const (
	// RevisionOrderingOlder lists newest revision first (default)
	RevisionOrderingOlder RevisionOrdering = "older"
	// RevisionOrderingNewer lists oldest revision first
	RevisionOrderingNewer RevisionOrdering = "newer"
)

type revisionsResponse struct {
	Batchcomplete bool                         `json:"batchcomplete"`
	Warnings      map[string]map[string]string `json:"warnings"`
	Query         struct {
		Normalized []struct {
			Fromencoded bool   `json:"fromencoded"`
			From        string `json:"from"`
			To          string `json:"to"`
		} `json:"normalized"`
		Pages []struct {
			PageID    int        `json:"pageid"`
			Ns        int        `json:"ns"`
			Title     string     `json:"title"`
			Revisions []Revision `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}
