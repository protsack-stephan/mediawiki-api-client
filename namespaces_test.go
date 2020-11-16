package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const namespacesTestURL = "/namespaces"
const namespacesTestID = 1
const namespacesTestName = "Test"

func createNamespacesServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(namespacesTestURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{ "batchcomplete": true, "query": { "namespaces": { "1": { "id": %d, "name": "%s" } } } }`, namespacesTestID, namespacesTestName)))
	})

	return router
}

func TestNamespaces(t *testing.T) {
	srv := httptest.NewServer(createNamespacesServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.NamespacesURL = namespacesTestURL

	ns, status, err := client.Namespaces(context.Background())

	if status != http.StatusOK {
		t.Fatal("namespaces response error")
	}

	if err != nil {
		t.Error(err)
	}

	if len(ns) != 1 {
		t.Fatal("namespaces count is wrong")
	}

	if ns[0].ID != namespacesTestID {
		t.Error("namespace id won't match")
	}

	if ns[0].Name != namespacesTestName {
		t.Error("namespace name won't match")
	}
}
