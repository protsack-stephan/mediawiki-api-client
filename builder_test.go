package mediawiki

import (
	"net/http"
	"testing"
	"time"
)

const builderTestURL = "http://localhost:5000"
const builderTestTimeout = time.Second * 2
const builderTestPageMetaURL = "/meta"
const builderTestPageHTMLURL = "/html"
const builderTestPageWikitextURL = "/wikitext"
const builderTestSitematrixURL = "/sitematrix"
const builderTestNamespacesURL = "/namespaces"
const builderTestPageRevisionsURL = "/revisions"

func TestBuilder(t *testing.T) {
	builder := NewBuilder(builderTestURL).
		HTTPClient(&http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
		}).
		Timeout(builderTestTimeout).
		Options(&Options{
			builderTestPageMetaURL,
			builderTestPageHTMLURL,
			builderTestPageWikitextURL,
			builderTestPageRevisionsURL,
			builderTestSitematrixURL,
			builderTestNamespacesURL,
		})

	client := builder.Build()

	if client.url != builderTestURL {
		t.Error("urls don't match")
	}

	if client.httpClient == nil {
		t.Error("http client is not set")
	}

	if client.options == nil {
		t.Fatal("http client options not set")
	}

	if client.options.PageMetaURL != builderTestPageMetaURL {
		t.Error("options meta url mismatch")
	}

	if client.options.PageHTMLURL != builderTestPageHTMLURL {
		t.Error("options html url mismatch")
	}

	if client.options.SitematrixURL != builderTestSitematrixURL {
		t.Error("sitematrix meta url mismatch")
	}

	if client.options.NamespacesURL != builderTestNamespacesURL {
		t.Error("namespace url mismatch")
	}

	if client.options.PageRevisionsURL != builderTestPageRevisionsURL {
		t.Error("namespace url mismatch")
	}

	if client.httpClient.Timeout != builderTestTimeout {
		t.Error("client timeout is not set")
	}

	if client.httpClient.Transport == nil {
		t.Error("client transport not set")
	}
}
