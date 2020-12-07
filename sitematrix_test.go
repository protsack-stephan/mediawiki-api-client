package mediawiki

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const sitematrixTestURL = "/sitematrix"
const sitematrixTestProjectName = "test_project"
const sitematrixTestProjectCode = "ua"
const sitematrixTestCount = 2

func createSitematrixServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(sitematrixTestURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf(`{ "sitematrix": { "count": %d, "0": { "code": "%s", "name": "%s",  "site": [] }, "specials": [] } }`, sitematrixTestCount, sitematrixTestProjectCode, sitematrixTestProjectName)))

		if err != nil {
			log.Panic(err)
		}
	})

	return router
}

func TestSitematrix(t *testing.T) {
	srv := httptest.NewServer(createSitematrixServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.SitematrixURL = sitematrixTestURL

	matrix, err := client.Sitematrix(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, sitematrixTestCount, matrix.Count)
	assert.Equal(t, 1, len(matrix.Projects))
	assert.Equal(t, sitematrixTestProjectName, matrix.Projects[0].Name)
	assert.Equal(t, sitematrixTestProjectCode, matrix.Projects[0].Code)
}
