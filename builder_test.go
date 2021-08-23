package mediawiki

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const builderTestURL = "http://localhost:5000"
const builderTestTimeout = time.Second * 2
const builderTestPageMetaURL = "/meta"
const builderTestPageHTMLURL = "/html"
const builderTestPageWikitextURL = "/wikitext"
const builderTestSitematrixURL = "/sitematrix"
const builderTestNamespacesURL = "/namespaces"
const builderTestPageRevisionsURL = "/revisions"
const builderTestPagesDataURL = "/pages-data"
const builderTestUserURL = "/users"

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
			builderTestPagesDataURL,
			builderTestUserURL,
		})

	client := builder.Build()

	assert.Equal(t, builderTestURL, client.url)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.options)
	assert.NotNil(t, client.httpClient.Transport)
	assert.Equal(t, builderTestPageMetaURL, client.options.PageMetaURL)
	assert.Equal(t, builderTestPageHTMLURL, client.options.PageHTMLURL)
	assert.Equal(t, builderTestSitematrixURL, client.options.SitematrixURL)
	assert.Equal(t, builderTestNamespacesURL, client.options.NamespacesURL)
	assert.Equal(t, builderTestPagesDataURL, client.options.PageDataURL)
	assert.Equal(t, builderTestUserURL, client.options.UserURL)
	assert.Equal(t, builderTestTimeout, client.httpClient.Timeout)
}
