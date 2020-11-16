package mediawiki

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const htmlTestURL = "/html"
const htmlTestTitle = "test_html"
const htmlTestRevision = 2
const htmlTestBody = "<h1>Hello world</h1>"

func createHTMLServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(htmlTestURL+htmlTestTitle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlTestBody))
	})

	router.HandleFunc(htmlTestURL+htmlTestTitle+"/"+strconv.Itoa(htmlTestRevision), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlTestBody))
	})

	return router
}

func TestPageHTML(t *testing.T) {
	srv := httptest.NewServer(createHTMLServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.PageHTMLURL = htmlTestURL

	html, status, err := client.PageHTML(context.Background(), htmlTestTitle)

	if status != http.StatusOK {
		t.Fatal("html response error")
	}

	if err != nil {
		t.Fatal(err)
	}

	if string(html) != htmlTestBody {
		t.Fatal("html responses don't match")
	}
}
