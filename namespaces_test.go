package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

	ns, err := client.Namespaces(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, 1, len(ns))
	assert.Equal(t, namespacesTestID, ns[0].ID)
	assert.Equal(t, namespacesTestName, ns[0].Name)
}
