package parser

import (
	"net/url"

	"golang.org/x/net/html"
)

var ImportantUrls = []string{
	"a",
	"link",
	"title",
	"meta",
	"h1",
	"article",
	"main",
	"img",
}

func ExtractLink(doc *html.Node , base *url.URL) {
	var links []string
	var walk func(node *html.Node)

	walk = func(node *html.Node) {
		for _ , li
	}
}
