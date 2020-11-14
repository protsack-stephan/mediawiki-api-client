package mediawiki

import "time"

const pageMetaURL = "/api/rest_v1/page/title/"

// PageMeta metadata for the page
type PageMeta struct {
	Title        string        `json:"title"`
	PageID       int           `json:"page_id"`
	Rev          int           `json:"rev"`
	Tid          string        `json:"tid"`
	Namespace    int           `json:"namespace"`
	UserID       int           `json:"user_id"`
	UserText     string        `json:"user_text"`
	Timestamp    time.Time     `json:"timestamp"`
	Comment      string        `json:"comment"`
	Tags         []string      `json:"tags"`
	Restrictions []interface{} `json:"restrictions"`
	PageLanguage string        `json:"page_language"`
	Redirect     bool          `json:"redirect"`
}

// PageMetaResponse response for page meta data
type PageMetaResponse struct {
	Items []PageMeta `json:"items"`
}
