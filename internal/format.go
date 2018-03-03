package internal

import (
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
)

func parseFromHTML(html string) ([]string, [][]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, nil, err
	}

	var header []string
	var datas [][]string
	var skipLine = -1
	doc.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		if skipLine == -1 && strings.Contains(tr.Find("td").Text(), "Edit Copy") {
			skipLine = 3
		}

		var data []string
		tr.Find("th").Each(func(_ int, th *goquery.Selection) {
			thText := th.Text()
			if thText != "" {
				header = append(header, thText)
			}
		})

		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			if i <= skipLine {
				return
			}
			data = append(data, td.Text())
		})

		if len(data) != 0 && (len(header) == 0 || (len(header) > 0 && len(header) == len(data))) {
			datas = append(datas, data)
		}
	})

	return header, datas, nil
}

func Format(name []string, values ...[]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(name)

	for _, v := range values {
		table.Append(v)
	}
	table.Render()
}

func FormatList(title string, values []string) {
	var vs [][]string
	for _, v := range values {
		vs = append(vs, []string{v})
	}
	Format([]string{title}, vs...)
}

// FromHTML Parse table from HTML
func ParseFromHTML(html string) error {
	header, datas, err := parseFromHTML(html)
	if err != nil {
		return err
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(header)
	for _, v := range datas {
		t.Append(v)
	}
	t.Render()

	return nil
}
