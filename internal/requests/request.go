package requests

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/Chyroc/phpmyadmin-cli/internal/common"
)

var c http.Client

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	c = http.Client{
		Jar: jar,
	}
}

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

	common.Debug("method %s\n", method)
	common.Debug("url %s\n", URL.String())
	common.Debug("body %#v\n", body)

	req, err := http.NewRequest(method, URL.String(), body)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	return c.Do(req)
}
