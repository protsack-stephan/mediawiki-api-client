package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pageRevisionsTestIDS = []int{944917628, 944912305}

const pageRevisionsTestURL = "/revisions/%d/%s"
const pageRevisionsTestTitle = "test"
const pageRevisionsTestLimit = 2

func createPageRevisionsServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf(pageRevisionsTestURL, pageRevisionsTestLimit, pageRevisionsTestTitle), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"continue":{"rvcontinue":"20200310174438|944912225","continue":"||"},"query":{"normalized":[{"fromencoded":false,"from":"Pet_door","to":"Pet door"}],"pages":[{"pageid":3276454,"ns":0,"title":"%s","revisions":[{"revid":%d,"parentid":944912305,"minor":true,"user":"IdreamofJeanie","timestamp":"2020-03-10T18:17:36Z","comment":"Reverted 2 edits by [[Special:Contributions\/112.196.52.43|112.196.52.43]] ([[User talk:112.196.52.43|talk]]) to last revision by LizzieBabes419 ([[WP:TW|TW]])"},{"revid":%d,"parentid":944912225,"minor":false,"user":"112.196.52.43","anon":true,"timestamp":"2020-03-10T17:45:04Z","comment":"\/* References *\/"}]}]}}`, pageRevisionsTestTitle, pageRevisionsTestIDS[0], pageRevisionsTestIDS[1])))
	})

	return router
}

func TestPageRevisions(t *testing.T) {
	srv := httptest.NewServer(createPageRevisionsServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.PageRevisionsURL = pageRevisionsTestURL

	revs, status, err := client.PageRevisions(context.Background(), pageRevisionsTestTitle, pageRevisionsTestLimit)

	assert.Equal(t, http.StatusOK, status)
	assert.Nil(t, err)
	assert.Equal(t, pageRevisionsTestLimit, len(revs))

	lookup := map[int]Revision{}

	for _, rev := range revs {
		lookup[rev.RevID] = rev
	}

	for _, revID := range pageRevisionsTestIDS {
		if _, ok := lookup[revID]; !ok {
			assert.True(t, ok)
		}
	}
}
