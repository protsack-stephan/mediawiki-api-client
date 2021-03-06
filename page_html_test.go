package mediawiki

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const htmlTestURL = "/html"
const htmlTestTitle = "test_html"
const htmlTestRevision = 2
const htmlTestBody = "<h1>Hello world</h1>"

func createHTMLServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(htmlTestURL+htmlTestTitle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(htmlTestBody))

		if err != nil {
			log.Panic(err)
		}
	})

	router.HandleFunc(htmlTestURL+htmlTestTitle+"/"+strconv.Itoa(htmlTestRevision), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(htmlTestBody))

		if err != nil {
			log.Panic(err)
		}
	})

	return router
}

func TestPageHTML(t *testing.T) {
	srv := httptest.NewServer(createHTMLServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.PageHTMLURL = htmlTestURL

	html, err := client.PageHTML(context.Background(), htmlTestTitle)

	assert.Nil(t, err)
	assert.Equal(t, htmlTestBody, string(html))
}
