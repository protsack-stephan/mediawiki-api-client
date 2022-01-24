package mediawiki

import "time"

const pageDataURL = "/w/api.php"

type PageDataOptions struct {
	RevisionsLimit int
	RevisionProps  []string
}

// PageDataOresScore representation for ORES score
type PageDataOresScore struct {
	True  float64 `json:"true"`
	False float64 `json:"false"`
}

// PageDataOresArticleQuality representation for ORES article quality
type PageDataOresArticleQuality struct {
	Stub float64 `json:"Stub"`
}

// PageDataOresScores representation of ORES scores
type PageDataOresScores struct {
	Articlequality PageDataOresArticleQuality `json:"articlequality"`
	Damaging       PageDataOresScore          `json:"damaging"`
	Goodfaith      PageDataOresScore          `json:"goodfaith"`
}

// PageDataRevisionMainSlot representation for revision main slot
type PageDataRevisionMainSlot struct {
	Contentmodel  string `json:"contentmodel"`
	Contentformat string `json:"contentformat"`
	Content       string `json:"content"`
}

// PageDataRevisionSlots representation for revision slots
type PageDataRevisionSlots struct {
	Main PageDataRevisionMainSlot `json:"main"`
}

// PageDataRevision revision data for pages data response
type PageDataRevision struct {
	RevID      int                   `json:"revid"`
	ParentID   int                   `json:"parentid"`
	User       string                `json:"user"`
	UserID     int                   `json:"userid"`
	Minor      bool                  `json:"minor"`
	Timestamp  time.Time             `json:"timestamp"`
	Slots      PageDataRevisionSlots `json:"slots"`
	Comment    string                `json:"comment"`
	Tags       []string              `json:"tags"`
	Oresscores interface{}           `json:"oresscores"`
}

// PageDataProtection representation for page data protection
type PageDataProtection struct {
	Type   string `json:"type"`
	Level  string `json:"level"`
	Expiry string `json:"expiry"`
}

// PageDataCategory representation for page data category
type PageDataCategory struct {
	Ns    int    `json:"ns"`
	Title string `json:"title"`
}

// PageDataTemplate representation for page data template
type PageDataTemplate struct {
	Ns    int    `json:"ns"`
	Title string `json:"title"`
}

// PageDataRedirect representation for page data redirect
type PageDataRedirect struct {
	PageID int    `json:"pageid"`
	Ns     int    `json:"ns"`
	Title  string `json:"title"`
}

// PageDataFlagged representation for page data flagged
type PageDataFlagged struct {
	StableRevID  int       `json:"stable_revid"`
	Level        int       `json:"level"`
	LevelText    string    `json:"level_text"`
	PendingSince time.Time `json:"pending_since"`
}

// PageDataProps representation for page data props
type PageDataProps struct {
	WikibaseItem string `json:"wikibase_item"`
}

// PageDataWebEntityUsage representation for page data web entity usage
type PageDataWebEntityUsage map[string]struct {
	Aspects []string `json:"aspects"`
}

// PageData page data returned from actions API
type PageData struct {
	PageID               int                    `json:"pageid"`
	Ns                   int                    `json:"ns"`
	Title                string                 `json:"title"`
	Watchers             int                    `json:"watchers"`
	ContentModel         string                 `json:"contentmodel"`
	PageLanguage         string                 `json:"pagelanguage"`
	PageLanguageHTMLCode string                 `json:"pagelanguagehtmlcode"`
	PageLanguageDir      string                 `json:"pagelanguagedir"`
	Touched              time.Time              `json:"touched"`
	LastRevID            int                    `json:"lastrevid"`
	Length               int                    `json:"length"`
	Missing              bool                   `json:"missing"`
	Protection           []PageDataProtection   `json:"protection"`
	Restrictiontypes     []string               `json:"restrictiontypes"`
	FullURL              string                 `json:"fullurl"`
	EditURL              string                 `json:"editurl"`
	CanonicalURL         string                 `json:"canonicalurl"`
	Displaytitle         string                 `json:"displaytitle"`
	Categories           []PageDataCategory     `json:"categories"`
	Revisions            []PageDataRevision     `json:"revisions"`
	Templates            []PageDataTemplate     `json:"templates"`
	Redirects            []PageDataRedirect     `json:"redirects"`
	Flagged              PageDataFlagged        `json:"flagged"`
	Pageprops            PageDataProps          `json:"pageprops"`
	WbEntityUsage        PageDataWebEntityUsage `json:"wbentityusage"`
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
