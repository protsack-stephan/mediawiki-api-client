package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const metaTestURL = "/meta"
const metaTestTitle = "test_title"
const metaTestRevision = 1

func createMetaServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(metaTestURL+metaTestTitle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"items": [ { "title": "%s", "rev": %d } ]}`, metaTestTitle, metaTestRevision)))
	})

	return router
}

func TestPageMeta(t *testing.T) {
	srv := httptest.NewServer(createMetaServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.PageMetaURL = metaTestURL

	meta, status, err := client.PageMeta(context.Background(), metaTestTitle)

	assert.Equal(t, http.StatusOK, status)
	assert.Nil(t, err)
	assert.Equal(t, metaTestTitle, meta.Title)
	assert.Equal(t, metaTestRevision, meta.Rev)
}
