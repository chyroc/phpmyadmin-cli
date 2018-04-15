package phpmyadmin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Chyroc/phpmyadmin-cli/internal/common"
	"github.com/Chyroc/phpmyadmin-cli/internal/requests"
)

var DefaultPHPMyAdmin *phpMyAdmin
var tokenRegexp = regexp.MustCompile("<input type=\"hidden\" name=\"token\" value=\"(.*?)\" [/]>")

type phpMyAdmin struct {
	*requests.Session
	Token string
	uri   string
}

type phpMyAdminResp struct {
	Message string
	Success bool
	Error   string
}

type Server struct {
	ID   string
	Name string
}

type Servers struct {
	S []Server
}

func handlerPhpmyadminResp(r phpMyAdminResp) ([]byte, error) {
	if !r.Success {
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

func (s *Servers) Print() {
	for _, v := range s.S {
		common.Info(fmt.Sprintf("%s: %s\n", v.ID, v.Name))
	}
}
func init() {
	DefaultPHPMyAdmin = &phpMyAdmin{
		Session: requests.DefaultSession,
	}
}

func (p *phpMyAdmin) SetURI(uri string) {
	p.uri = uri
}

func (p *phpMyAdmin) requestGET(uri string) ([]byte, error) {
	resp, err := p.Get(p.uri, uri, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (p *phpMyAdmin) initCookie() error {
	resp, err := requests.DefaultSession.Get(p.uri, "index.php", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	matchToken := tokenRegexp.FindStringSubmatch(string(b))
	if len(matchToken) != 2 {
		return fmt.Errorf("match token err: %s", strings.Join(matchToken, ";"))
	} else if matchToken[1] == "" {
		return fmt.Errorf("empty token")
	}

	p.Token = matchToken[1]

	return nil
}

func (p *phpMyAdmin) Login(username, password string) (err error) {
	defer func() {
		if err != nil {
			common.Error(err)
		}
	}()

	body := strings.NewReader(fmt.Sprintf("pma_username=%s&pma_password=%s&lang=en", username, password))
	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp, err := requests.DefaultSession.Post(p.uri, "index.php", nil, header, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("\n\n\nindex.php %s\n\n\n", b)

	result, err := p.requestGET("server_status_processes.php")
	if err != nil {
		return err
	}
	fmt.Printf("\n\n\nstatus result %s\n\n\n", result)
	if !strings.Contains(string(result), "SHOW PROCESSLIST") {
		return fmt.Errorf("login err")
	}

	common.Info("login as [%s] success\n", username)
	return nil
}

func (p *phpMyAdmin) GetServerList(url string) (*Servers, error) {
	resp, err := p.Get(url, "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if common.IsDebug1 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			common.Debug("err %s\n", err)
		}
		common.Debug("return %s\n", string(b))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var s []Server
	doc.Find("#select_server > option").Each(func(_ int, selection *goquery.Selection) {
		id := strings.TrimSpace(selection.AttrOr("value", ""))
		name := strings.TrimSpace(selection.Text())

		if id != "" {
			s = append(s, Server{id, name})
		}
	})

	return &Servers{s}, nil
}

func (p *phpMyAdmin) ExecSQL(server, database, table, sql string) ([]byte, error) {
	if p.Token == "" {
		p.initCookie()
	}

	data := map[string]string{
		// "table":             table,
		"db":                database,
		"server":            server,
		"token":             p.Token,
		"prev_sql_query":    "",
		"sql_query":         sql,
		"ajax_request":      "true",
		"ajax_page_request": "true",
	}
	var bs []string
	for k, v := range data {
		bs = append(bs, k+"="+url.QueryEscape(v))
	}
	body := strings.NewReader(strings.Join(bs, "&"))
	common.Debug("ExecSQL [%v]\n", strings.Join(bs, "&"))
	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

	resp, err := p.Post(p.uri+"/import.php", "", nil, header, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r phpMyAdminResp
	if err = json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	common.Debug("ExecSQL [%v]:[%v]:[%v]\n", r.Success, r.Error, r.Message)

	return handlerPhpmyadminResp(r)
}
