package requests

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func request(method, uri, path string, query, header map[string]string, body io.Reader, cookies []*http.Cookie) (*http.Response, error) {
	// todo fix and test http://http//xxx
	if !strings.HasPrefix(uri, "http://") || !strings.HasPrefix(uri, "https://") {
		uri = "http://" + uri
	}

	URL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if path != "" {
		URL.Path = path
	}

	if query != nil {
		q := URL.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		URL.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, URL.String(), body)
	if err != nil {
		return nil, err
	}

	for _, v := range cookies {
		req.AddCookie(v)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	return http.DefaultClient.Do(req)
}
