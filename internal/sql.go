package internal

import (
	"strings"

	"fmt"
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
	var fields [][]string
	var nFields []string
	body.Find("th").Each(func(i int, selection *goquery.Selection) {
		if selection != nil {
			//field := strings.Split(selection.Text(), "\n")
			//fmt.Printf("", field)
			//fields = append(fields, field)
			//nFields = append(nFields, field[0])
			if t := selection.Text(); t != "" {
				nFields = append(nFields, t)
			}
		}
	})

	fmt.Printf("", fields)
	fmt.Printf("", nFields)

	var values [][]string
	body.Find("tr").Each(func(i int, selection *goquery.Selection) {
		if selection != nil {
			var value []string
			body.Find("td").Each(func(i int, selection *goquery.Selection) {
				if selection != nil {
					value = append(value, selection.Text())
				}
			})
			values = append(values, value)
		}
	})

	var nValues [][]string
	for _, v := range values {
		nValues = append(nValues, v[len(v)-len(nFields):])
	}
	//fmt.Printf("", nValues)

	return nFields, nValues
}
