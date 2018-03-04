package internal

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func filterEmpty(vs []string) []string {
	var vss []string
	for _, v := range vs {
		if v != "" {
			vss = append(vss, v)
		}
	}

	return vss
}
func ToList(selecttion *goquery.Selection) []string {
	s := strings.TrimSpace(selecttion.Text())
	return filterEmpty(strings.Split(s, "\n"))
}
