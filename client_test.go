package mediawiki

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const metaTestTitle = "test_title"
const metaTestRevision = 1
const htmlTestTitle = "test_html"
const htmlTestRevision = 2
const htmlTestBody = "<h1>Hello world</h1>"

func createClientServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(pageHTMLURL+htmlTestTitle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlTestBody))
	})

	router.HandleFunc(pageHTMLURL+htmlTestTitle+"/"+strconv.Itoa(htmlTestRevision), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlTestBody))
	})

	router.HandleFunc(pageMetaURL+metaTestTitle, func(w http.ResponseWriter, r *http.Request) {
		meta := new(PageMetaResponse)
		meta.Items = []PageMeta{
			{
				Title: metaTestTitle,
				Rev:   metaTestRevision,
			},
		}

		body, err := json.Marshal(meta)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	return router
}

func TestClient(t *testing.T) {
	srv := httptest.NewServer(createClientServer())
	defer srv.Close()

	client := NewClient(srv.URL)

	if client == nil {
		t.Fatal("client is empty")
	}

	meta, status, err := client.PageMeta(context.Background(), metaTestTitle)

	if status != http.StatusOK {
		t.Fatal("meta response error")
	}

	if err != nil {
		t.Fatal(err)
	}

	if meta.Title != metaTestTitle {
		t.Error("meta titles don't mach")
	}

	if meta.Rev != metaTestRevision {
		t.Error("meta revisions don't mach")
	}

	html, status, err := client.PageHTML(context.Background(), htmlTestTitle)

	if status != http.StatusOK {
		t.Fatal("html response error")
	}

	if err != nil {
		t.Fatal(err)
	}

	if string(html) != htmlTestBody {
		t.Fatal("html response don't mach")
	}
}
