package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pageRevisionsTestIDS = []int{944917628, 944912305}

const (
	pageRevisionsTestTitle = "test"
	pageRevisionsTestLimit = 2
	pageRevisionsTestOrder = RevisionOrderingNewer
	pageRevisionsTestProps = "content|ids"
)

func createPageRevisionsServer(t *testing.T) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		values, err := url.ParseQuery(r.URL.RawQuery)
		assert.NoError(t, err)

		assert.True(t, values.Has("titles") && assert.Equal(t, pageRevisionsTestTitle, values.Get("titles")))
		assert.True(t, values.Has("rvlimit") && assert.Equal(t, strconv.Itoa(pageRevisionsTestLimit), values.Get("rvlimit")))
		assert.True(t, values.Has("rvdir") && assert.Equal(t, string(pageRevisionsTestOrder), values.Get("rvdir")))
		assert.True(t, values.Has("rvprop") && assert.Equal(t, pageRevisionsTestProps, values.Get("rvprop")))

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(fmt.Sprintf(`{"continue":{"rvcontinue":"20200310174438|944912225","continue":"||"},"query":{"normalized":[{"fromencoded":false,"from":"Pet_door","to":"Pet door"}],"pages":[{"pageid":3276454,"ns":0,"title":"%s","revisions":[{"revid":%d,"parentid":944912305,"minor":true,"user":"IdreamofJeanie","timestamp":"2020-03-10T18:17:36Z","comment":"Reverted 2 edits by [[Special:Contributions\/112.196.52.43|112.196.52.43]] ([[User talk:112.196.52.43|talk]]) to last revision by LizzieBabes419 ([[WP:TW|TW]])"},{"revid":%d,"parentid":944912225,"minor":false,"user":"112.196.52.43","anon":true,"timestamp":"2020-03-10T17:45:04Z","comment":"\/* References *\/"}]}]}}`, pageRevisionsTestTitle, pageRevisionsTestIDS[0], pageRevisionsTestIDS[1])))

		assert.NoError(t, err)
	})

	return router
}

func TestPageRevisions(t *testing.T) {
	srv := httptest.NewServer(createPageRevisionsServer(t))
	defer srv.Close()

	client := NewClient(srv.URL)

	revs, err := client.PageRevisions(
		context.Background(),
		pageRevisionsTestTitle,
		pageRevisionsTestLimit,
		PageRevisionsOptions{
			Order: pageRevisionsTestOrder,
			Props: []string{"content", "ids"},
		},
	)

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
