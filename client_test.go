package mediawiki

import (
	"testing"
)

const clientTestURL = "http://localhost:5000/"

func TestClient(t *testing.T) {
	client := NewClient(clientTestURL)

	if client == nil {
		t.Fatal("client is empty")
	}

	if client.url != clientTestURL {
		t.Error("client url is wrong")
	}
}
