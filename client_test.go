package mediawiki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const clientTestURL = "http://localhost:5000/"

func TestClient(t *testing.T) {
	client := NewClient(clientTestURL)
	assert.NotNil(t, client)
	assert.Equal(t, clientTestURL, client.url)
}
