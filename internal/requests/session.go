package requests

import (
	"io"
	"net/http"
	"sync"
	"net/http/cookiejar"
)

type Session struct {
	cookies    map[string]*http.Cookie
	cookieLock *sync.Mutex
	client     http.Client
}

var DefaultSession *Session

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	DefaultSession = &Session{
		cookies:    make(map[string]*http.Cookie),
		cookieLock: new(sync.Mutex),
		client: http.Client{
			Jar: jar,
		},
	}
}

func (r *Session) setCookie(resp *http.Response) {
	r.cookieLock.Lock()
	defer r.cookieLock.Unlock()

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		r.cookies[cookie.Name] = cookie
	}
}

func (r *Session) readCookie() []*http.Cookie {
	r.cookieLock.Lock()
	defer r.cookieLock.Unlock()

	var cookies []*http.Cookie
	for _, cookie := range r.cookies {
		cookies = append(cookies, cookie)
	}

	return cookies
}

func (r *Session) Cookie(k string) string {
	c, ok := r.cookies[k]
	if !ok {
		return ""
	}

	return c.Value
}

func (r *Session) Get(uri, path string, query map[string]string) (*http.Response, error) {
	resp, err := request(r, http.MethodGet, uri, path, query, nil, nil)
	if err != nil {
		return nil, err
	}

	r.setCookie(resp)

	return resp, nil
}

func (r *Session) Post(uri, path string, query, header map[string]string, body io.Reader) (*http.Response, error) {
	resp, err := request(r, http.MethodPost, uri, path, query, header, body)
	if err != nil {
		return nil, err
	}

	r.setCookie(resp)

	return resp, nil
}
