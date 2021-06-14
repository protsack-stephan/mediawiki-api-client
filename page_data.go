package mediawiki

import "time"

const pageDataURL = "/w/api.php"

// PageData page data returned from actions API
type PageData struct {
	PageID               int       `json:"pageid"`
	Ns                   int       `json:"ns"`
	Title                string    `json:"title"`
	ContentModel         string    `json:"contentmodel"`
	PageLanguage         string    `json:"pagelanguage"`
	PageLanguageHTMLCode string    `json:"pagelanguagehtmlcode"`
	PageLanguageDir      string    `json:"pagelanguagedir"`
	Touched              time.Time `json:"touched"`
	LastRevID            int       `json:"lastrevid"`
	Length               int       `json:"length"`
	Missing              bool      `json:"missing"`
	Protection           []struct {
		Type   string `json:"type"`
		Level  string `json:"level"`
		Expiry string `json:"expiry"`
	} `json:"protection"`
	Restrictiontypes []string `json:"restrictiontypes"`
	FullURL          string   `json:"fullurl"`
	EditURL          string   `json:"editurl"`
	CanonicalURL     string   `json:"canonicalurl"`
	Displaytitle     string   `json:"displaytitle"`
	Categories       []struct {
		Ns    int    `json:"ns"`
		Title string `json:"title"`
	} `json:"categories"`
	Revisions []struct {
		RevID     int       `json:"revid"`
		ParentID  int       `json:"parentid"`
		User      string    `json:"user"`
		Timestamp time.Time `json:"timestamp"`
		Slots     struct {
			Main struct {
				Contentmodel  string `json:"contentmodel"`
				Contentformat string `json:"contentformat"`
				Content       string `json:"content"`
			} `json:"main"`
		} `json:"slots"`
		Comment    string   `json:"comment"`
		Tags       []string `json:"tags"`
		Oresscores struct {
			Articlequality struct {
				Stub float64 `json:"Stub"`
			} `json:"articlequality"`
			Damaging struct {
				True  float64 `json:"true"`
				False float64 `json:"false"`
			} `json:"damaging"`
			Goodfaith struct {
				True  float64 `json:"true"`
				False float64 `json:"false"`
			} `json:"goodfaith"`
		} `json:"oresscores"`
	} `json:"revisions"`
	Templates []struct {
		Ns    int    `json:"ns"`
		Title string `json:"title"`
	} `json:"templates"`
	WbEntityUsage map[string]struct {
		Aspects []string `json:"aspects"`
	} `json:"wbentityusage"`
	Pageprops struct {
		WikibaseItem string `json:"wikibase_item"`
	} `json:"pageprops"`
	Redirects []struct {
		PageID int    `json:"pageid"`
		Ns     int    `json:"ns"`
		Title  string `json:"title"`
	} `json:"redirects"`
}

type pageDataResponse struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Normalized []struct {
			Fromencoded bool   `json:"fromencoded"`
			From        string `json:"from"`
			To          string `json:"to"`
		} `json:"normalized"`
		Redirects []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"redirects"`
		Pages []PageData `json:"pages"`
	} `json:"query"`
}
