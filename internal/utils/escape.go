package utils

import (
	"strings"
	"net/url"
)

func Escape(s string) string {
	s = strings.Replace(s, "&amp;", "&", -1)
	s = strings.Replace(s, "&gt;", ">", -1)
	s = strings.Replace(s, "&lt;", "<", -1)
	s = strings.Replace(s, "&quot;", "\"", -1)

	return url.QueryEscape(s)
}
