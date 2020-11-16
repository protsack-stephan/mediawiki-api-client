package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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

	if status != http.StatusOK {
		t.Fatal("meta response error")
	}

	if err != nil {
		t.Fatal(err)
	}

	if meta.Title != metaTestTitle {
		t.Error("meta titles don't match")
	}

	if meta.Rev != metaTestRevision {
		t.Error("meta revisions don't match")
	}
}
