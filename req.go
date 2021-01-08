package mediawiki

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
)

func req(ctx context.Context, cl *http.Client, method string, url string, reqBody io.Reader, headers ...map[string]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)

	for _, header := range headers {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}

	if err != nil {
		return nil, 0, err
	}

	res, err := cl.Do(req)

	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, res.StatusCode, err
	}

	return resBody, res.StatusCode, nil
}
