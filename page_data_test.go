package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const pageDataTestURL = "/page-data"
const pageDataTestTitle = "Ninja"
const pageDataTestQID = "Q90"
const pageDataTestRev = 998092778
const pageDataTestBdy = `{
	"batchcomplete": true,
	"query": {
			"pages": [
					{
							"pageid": 22989,
							"ns": 0,
							"title": "%s",
							"pageprops": {
									"wikibase_item": "%s"
							},
							"contentmodel": "wikitext",
							"pagelanguage": "en",
							"pagelanguagehtmlcode": "en",
							"pagelanguagedir": "ltr",
							"touched": "2021-01-05T07:56:19Z",
							"lastrevid": %d,
							"length": 263051,
							"revisions": [
									{
											"revid": 998092778,
											"parentid": 998092404,
											"minor": false,
											"user": "Politicsfan4",
											"timestamp": "2021-01-03T19:49:57Z",
											"comment": "Reverted 1 pending edit by [[Special:Contributions/174.246.137.154|174.246.137.154]] to revision 997918021 by Maxeto0910"
									}
							]
					}
			]
	}
}`

func createPageDataServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(pageDataTestURL, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(pageDataTestBdy, pageDataTestTitle, pageDataTestQID, pageDataTestRev)))
	})

	return router
}

func TestPageData(t *testing.T) {
	srv := httptest.NewServer(createPageDataServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.PageDataURL = fmt.Sprintf("%s%s", pageDataTestURL, "?title=%s")

	pages, err := client.PagesData(context.Background(), pageDataTestTitle)
	assert.NoError(t, err)

	for _, page := range pages {
		assert.Equal(t, pageDataTestTitle, page.Title)
		assert.Equal(t, pageDataTestQID, page.Pageprops.WikibaseItem)
		assert.Equal(t, pageDataTestRev, page.LastRevID)
	}
}
