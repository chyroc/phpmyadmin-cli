package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PhpmyadminResp struct {
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

func (s *Servers) Print() {
	for _, v := range s.S {
		Info(fmt.Sprintf("%s: %s\n", v.ID, v.Name))
	}
}

func GetServerList(url string) (*Servers, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s", url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
