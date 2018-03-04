package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Request(uri, db, sql string) (*goquery.Selection, error) {
	var url string
	if db == "" {
		url = fmt.Sprintf("http://%s/import.php?sql_query=%s&show_query=1&server=3&ajax_request=true", uri, sql)
	} else {
		url = fmt.Sprintf("http://%s/import.php?db=%s&sql_query=%s&show_query=1&server=3&ajax_request=true", uri, db, sql)
	}

	//fmt.Printf("url: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r struct {
		Message string
		Success bool
		Error   string
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	} else if !r.Success {
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

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r.Message))
	if err != nil {
		return nil, err
	}

	return doc.Find(".table_results"), nil
}
