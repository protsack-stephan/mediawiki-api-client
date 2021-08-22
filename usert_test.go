package mediawiki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const userTestURL = "/user"
const userTestID = 100
const userTestMissingID = 999
const userTestName = "Ninja"
const userTestEditCount = 2
const userTestBody = `{
	"batchcomplete": true,
	"query": {
			"users": [
					{
							"userid": %d,
							"name": "%s",
							"editcount": %d,
							"registration": "2021-04-02T13:43:05Z",
							"groups": [
									"*",
									"user",
									"autoconfirmed"
							],
							"groupmemberships": [],
							"emailable": false
					},
					{
							"userid": %d,
							"missing": true
					}
			]
	}
}`

func createUserServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(userTestURL, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(
			userTestBody,
			userTestID,
			userTestName,
			userTestEditCount,
			userTestMissingID)))
	})

	return router
}

func assertUser(assert *assert.Assertions, user User) {
	assert.Equal(userTestID, user.UserID)
	assert.Equal(userTestName, user.Name)
	assert.Equal(userTestEditCount, user.EditCount)
	assert.NotEmpty(user.Groups)
}

func TestUser(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	srv := httptest.NewServer(createUserServer())
	defer srv.Close()

	client := NewClient(srv.URL)
	client.options.UserURL = userTestURL

	users, err := client.Users(ctx, userTestID, userTestMissingID)
	assert.NoError(err)
	assert.Contains(users, userTestID)
	assert.NotContains(users, userTestMissingID)

	for id, user := range users {
		assert.Equal(userTestID, id)
		assertUser(assert, user)
	}

	_, err = client.User(ctx, userTestMissingID)
	assert.Equal(ErrUserNotFound, err)

	user, err := client.User(ctx, userTestID)
	assert.NoError(err)
	assertUser(assert, user)
}
