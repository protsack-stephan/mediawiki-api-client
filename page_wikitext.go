package mediawiki

const pageWikitextURL = "/w/api.php?action=query&format=json&prop=revisions&formatversion=2&titles=%s&rvprop=content&rvslots=main&rvlimit=1"

type wikitextResponse struct {
	Continue struct {
		RvContinue string `json:"rvcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Pages []struct {
			PageID    int    `json:"pageid"`
			Ns        int    `json:"ns"`
			Title     string `json:"title"`
			Revisions []struct {
				Slots struct {
					Main struct {
						Contentmodel  string `json:"contentmodel"`
						Contentformat string `json:"contentformat"`
						Content       string `json:"content"`
					} `json:"main"`
				} `json:"slots"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}
