package mediawiki

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const reqTestURL = "/test"
const reqTestBody = "test"

func createReqServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(reqTestURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(reqTestBody))

		if err != nil {
			log.Panic(err)
		}
	})

	return router
}

func TestReq(t *testing.T) {
	srv := httptest.NewServer(createReqServer())
	defer srv.Close()

	body, status, err := req(context.Background(), http.DefaultClient, http.MethodGet, srv.URL+reqTestURL, nil)

	assert.Equal(t, http.StatusOK, status)
	assert.Nil(t, err)
	assert.Equal(t, reqTestBody, string(body))
}
