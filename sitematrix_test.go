package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const sitematrixTestURL = "/sitematrix"
const sitematrixProjectName = "test_project"
const sitematrixProjectCode = "ua"
const sitematrixCount = 2

func createSitematrixServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(sitematrixTestURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{ "sitematrix": { "count": %d, "0": { "code": "%s", "name": "%s",  "site": [] }, "specials": [] } }`, sitematrixCount, sitematrixProjectCode, sitematrixProjectName)))
	})

	return router
}

func TestSitematrix(t *testing.T) {
	srv := httptest.NewServer(createSitematrixServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.SitematrixURL = sitematrixTestURL

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
