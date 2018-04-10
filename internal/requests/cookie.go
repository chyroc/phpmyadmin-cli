package requests

import (
	"io"
	"net/http"
	"sync"
)

type Session struct {
	cookies    map[string]*http.Cookie
	cookieLock *sync.Mutex
}

var DefaultSession *Session

func init() {
	DefaultSession = &Session{
		cookies:    make(map[string]*http.Cookie),
		cookieLock: new(sync.Mutex),
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
	resp, err := request(http.MethodGet, uri, path, query, nil, nil, r.readCookie())
	if err != nil {
		return nil, err
	}

	r.setCookie(resp)

	return resp, nil
}

func (r *Session) Post(uri, path string, query, header map[string]string, body io.Reader) (*http.Response, error) {
	resp, err := request(http.MethodPost, uri, path, query, header, body, r.readCookie())
	if err != nil {
		return nil, err
	}

	r.setCookie(resp)

	return resp, nil
}
