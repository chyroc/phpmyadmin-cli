package phpmyadmin

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func docDatabases(resp *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	var databases []string
	doc.Find("#pma_navigation_tree_content > ul > li").Each(func(i int, selection *goquery.Selection) {
		d := strings.TrimSpace(selection.Find("li > a").Text())
		if d != "" {
			databases = append(databases, d)
		}
	})

	return databases, nil
}
