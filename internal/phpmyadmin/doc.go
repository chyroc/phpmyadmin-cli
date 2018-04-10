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

func docTables(text string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	var tables []string
	doc.Find("tbody > tr").Each(func(i int, selection *goquery.Selection) {
		d := strings.TrimSpace(selection.Find("tr > th > a").Text())
		if d != "" {
			tables = append(tables, d)
		}

	})

	return tables, nil
}
