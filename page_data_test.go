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
const pageDataTestRedirectTitle = "Redirect"
const pageDataTestMissingTitle = "NotThere"
const pageDataTestQID = "Q90"
const pageDataTestRev = 998092778
const pageDataTestWikitext = "...wikitext goes here..."
const pageDataTestRedirectsPageID = 903711
const pageDataTestRedirectsTitle = "Super Ninja"
const pageDataTestRedirectsNs = 2
const pageDataTestTemplatesTitle = "Template:Katana"
const pageDataTestTemplatesNs = 10
const pageDataTestWbEntityUsageQID = "P569"
const pageDataTestWbEntityUsageAspect = "O"
const pageDataTestCategoriesTitle = "Category:Japan"
const pageDataTestCategoriesNs = 14
const pageDataTestBdy = `{
	"batchcomplete": true,
	"query": {
			"redirects": [
					{
							"from": "Redirect",
							"to": "Redirect-1"
					}
			],
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
							"redirects": [
									{
											"pageid": %d,
											"ns": %d,
											"title": "%s"
									}
							],
							"templates": [
									{
											"ns": %d,
											"title": "%s"
									}
							],
							"wbentityusage": {
									"%s": {
											"aspects": [
													"%s",
													"D.en",
													"O",
													"S",
													"T"
											]
									}
								},
							"categories": [
								{
									"ns": %d,
									"title": "%s"
								}
							],
							"revisions": [
									{
											"revid": 998092778,
											"parentid": 998092404,
											"minor": false,
											"user": "Politicsfan4",
											"timestamp": "2021-01-03T19:49:57Z",
											"comment": "Reverted 1 pending edit by [[Special:Contributions]] to revision 997918021",
											"slots": {
												"main": {
														"contentmodel": "wikitext",
														"contentformat": "text/x-wiki",
														"content": "%s"
												}
											}
									}
							]
					},
					{
						"pageid": 22989,
						"ns": 0,
						"title": "Redirect-1",
						"contentmodel": "wikitext",
						"pagelanguage": "en",
						"pagelanguagehtmlcode": "en",
						"pagelanguagedir": "ltr",
						"touched": "2021-01-05T07:56:19Z",
						"lastrevid": 100,
						"length": 263051,
						"revisions": [
								{
										"revid": 998092778,
										"parentid": 998092404,
										"minor": false,
										"user": "Politicsfan4",
										"timestamp": "2021-01-03T19:49:57Z",
										"comment": "Reverted 1 pending edit by [[Special:Contributions]] to revision 997918021"
								}
						]
					},
					{
						"ns": 0,
						"title": "%s",
						"missing": true,
						"contentmodel": "wikitext",
						"pagelanguage": "af",
						"pagelanguagehtmlcode": "af",
						"pagelanguagedir": "ltr"
					}
			]
	}
}`

func createPageDataServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(pageDataTestURL, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(
			pageDataTestBdy,
			pageDataTestTitle,
			pageDataTestQID,
			pageDataTestRev,
			pageDataTestRedirectsPageID,
			pageDataTestRedirectsNs,
			pageDataTestRedirectsTitle,
			pageDataTestTemplatesNs,
			pageDataTestTemplatesTitle,
			pageDataTestWbEntityUsageQID,
			pageDataTestWbEntityUsageAspect,
			pageDataTestCategoriesNs,
			pageDataTestCategoriesTitle,
			pageDataTestWikitext,
			pageDataTestMissingTitle)))
	})

	return router
}

func TestPageData(t *testing.T) {
	assert := assert.New(t)
	srv := httptest.NewServer(createPageDataServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.PageDataURL = pageDataTestURL

	pages, err := client.PagesData(context.Background(), pageDataTestTitle, pageDataTestRedirectTitle, pageDataTestMissingTitle)
	assert.NoError(err)
	assert.NotContains(pages, pageDataTestMissingTitle)
	assert.NotContains(pages, pageDataTestRedirectTitle)

	for title, page := range pages {
		assert.Equal(pageDataTestTitle, title)
		assert.Equal(pageDataTestQID, page.Pageprops.WikibaseItem)
		assert.Equal(pageDataTestRev, page.LastRevID)
		assert.Equal(pageDataTestWikitext, page.Revisions[0].Slots.Main.Content)
		assert.Equal(pageDataTestRedirectsPageID, page.Redirects[0].PageID)
		assert.Equal(pageDataTestRedirectsNs, page.Redirects[0].Ns)
		assert.Equal(pageDataTestRedirectsTitle, page.Redirects[0].Title)
		assert.Equal(pageDataTestTemplatesTitle, page.Templates[0].Title)
		assert.Equal(pageDataTestTemplatesNs, page.Templates[0].Ns)
		assert.Contains(page.WbEntityUsage, pageDataTestWbEntityUsageQID)
		assert.NotEmpty(page.WbEntityUsage[pageDataTestWbEntityUsageQID].Aspects)
		assert.Contains(page.WbEntityUsage[pageDataTestWbEntityUsageQID].Aspects, pageDataTestWbEntityUsageAspect)
		assert.Equal(pageDataTestCategoriesTitle, page.Categories[0].Title)
		assert.Equal(pageDataTestCategoriesNs, page.Categories[0].Ns)
	}
}
