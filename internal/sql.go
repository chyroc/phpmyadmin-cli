package internal

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"fmt"
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
	a := filterEmpty(strings.Split(s, "\n"))
	fmt.Printf("%d %v\n", len(a), a)
	return a
}
