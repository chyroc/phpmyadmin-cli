package phpmyadmin

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Chyroc/phpmyadmin-cli/internal/common"
)

func handlerPhpmyadminResp(r phpMyAdminResp) ([]byte, error) {
	if r.Message == "" && r.Error == "" && !r.Success {
		return nil, fmt.Errorf("invalid PhpmyadminResp")
	}

	if !r.Success {
		common.Debug("%#v\n", r)
		errdoc, err := goquery.NewDocumentFromReader(strings.NewReader(r.Error))
		if err != nil {
			return nil, err
		}

		var errMessage string
		errdoc.Find("code").Each(func(_ int, code *goquery.Selection) {
			if code.Text() != "" {
				errMessage = code.Text()
			}
		})
		return nil, fmt.Errorf(errMessage)
	}

	return []byte(r.Message), nil
}

func refreshToken(p *phpMyAdmin, resp *http.Response) ([]byte, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	matchToken := tokenRegexp.FindStringSubmatch(string(b))
	if len(matchToken) == 2 {
		p.Token = matchToken[1]
	}

	return b, nil
}

func (p *phpMyAdmin) Get(uri, path string, query map[string]string) ([]byte, error) {
	resp, err := p.session.Get(uri, path, query)
	if err != nil {
		return nil, err
	}

	return refreshToken(p, resp)
}

func (p *phpMyAdmin) Post(uri, path string, query, header map[string]string, body io.Reader) ([]byte, error) {
	resp, err := p.session.Post(uri, path, query, header, body)
	if err != nil {
		return nil, err
	}

	return refreshToken(p, resp)
}
