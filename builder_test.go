package mediawiki

import (
	"net/http"
	"testing"
	"time"
)

func TestBuilder(t *testing.T) {
	const url = "http://test.com"
	const timeout = time.Second * 2
	const pageMetaURL = "/meta"
	const pageHTMLURL = "/html"

	builder := NewBuilder(url).
		HTTPClient(&http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
		}).
		Timeout(timeout).
		Options(&Options{
			PageMetaURL: pageMetaURL,
			PageHTMLURL: pageHTMLURL,
		})

	client := builder.Build()

	if client.url != url {
		t.Error("urls don't match")
	}

	if client.httpClient == nil {
		t.Error("http client is not set")
	}

	if client.options == nil {
		t.Fatal("http client options not set")
	}

	if client.options.PageHTMLURL != pageHTMLURL {
		t.Error("options html url mismatch")
	}

	if client.options.PageMetaURL != pageMetaURL {
		t.Error("options meta url mismatch")
	}

	if client.httpClient.Timeout != timeout {
		t.Error("client timeout is not set")
	}

	if client.httpClient.Transport == nil {
		t.Error("client transport not set")
	}
}
