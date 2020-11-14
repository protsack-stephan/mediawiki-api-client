package mediawiki

import (
	"context"
	"fmt"
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

const sitematrixTestURL = "/sitematrix"
const sitematrixProjectName = "test_project"
const sitematrixProjectCode = "ua"
const sitematrixCount = 2

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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"items": [ { "title": "%s", "rev": %d } ]}`, metaTestTitle, metaTestRevision)))
	})

	router.HandleFunc(sitematrixTestURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{ "sitematrix": { "count": %d, "0": { "code": "%s", "name": "%s",  "site": [] }, "specials": [] } }`, sitematrixCount, sitematrixProjectCode, sitematrixProjectName)))
	})

	return router
}

func TestClient(t *testing.T) {
	srv := httptest.NewServer(createClientServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.SitematrixURL = sitematrixTestURL

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
		t.Error("meta titles don't match")
	}

	if meta.Rev != metaTestRevision {
		t.Error("meta revisionses don't match")
	}

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

	matrix, status, err := client.Sitematrix(context.Background())

	if status != http.StatusOK {
		t.Fatal("matrix response error")
	}

	if err != nil {
		t.Fatal(err)
	}

	if matrix.Count != sitematrixCount {
		t.Error("matrix counts don't match")
	}

	if len(matrix.Projects) != 1 {
		t.Fatal("matrix projects count is wrong")
	}

	if matrix.Projects[0].Name != sitematrixProjectName {
		t.Error("matrix project names don't match")
	}

	if matrix.Projects[0].Code != sitematrixProjectCode {
		t.Error("matrix project codes don't mach")
	}
}
