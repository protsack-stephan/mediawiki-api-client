package mediawiki

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const reqTestURl = "/test"
const reqTestBody = "test"

func createReqServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(reqTestURl, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(reqTestBody))
	})

	return router
}

func TestReq(t *testing.T) {
	srv := httptest.NewServer(createReqServer())
	defer srv.Close()

	body, status, err := req(context.Background(), http.DefaultClient, http.MethodGet, srv.URL+reqTestURl, nil)

	assert.Equal(t, http.StatusOK, status)
	assert.Nil(t, err)
	assert.Equal(t, reqTestBody, string(body))
}
