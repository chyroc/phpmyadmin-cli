package internal

import (
	"fmt"
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
	a := filterEmpty(strings.Split(s, "\n"))
	fmt.Printf("%d %v\n", len(a), a)
	return a
}

func ToSelectData(body *goquery.Selection) ([]string, [][]string) {
	var nFields []string
	body.Find("th").Each(func(i int, selection *goquery.Selection) {
		if selection != nil {
			if t := selection.Text(); t != "" {
				nFields = append(nFields, t)
			}
		}
	})

	var values [][]string
	body.Find("tr").Each(func(i int, tr *goquery.Selection) {
		var tds []string
		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			tds = append(tds, td.Text())
		})
		if len(tds) != 0 {
			values = append(values, tds)
		}
	})

	var value []string
	body.Find("td").Each(func(i int, selection *goquery.Selection) {
		if selection != nil {
			value = append(value, selection.Text())
		}
	})

	var nValues [][]string
	for _, v := range values {
		if len(v)-len(nFields) >= len(v) {
			continue
		}
		fmt.Printf("%d %s", len(v), v)
		nValues = append(nValues, v[len(v)-len(nFields):])
	}

	return nFields, nValues
}
