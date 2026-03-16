package yjsy

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func safeExtractHTMLFirst(node *html.Node, expr string) string {
	res := htmlquery.FindOne(node, expr)

	if res == nil {
		return ""
	}

	return htmlquery.OutputHTML(res, false)
}
