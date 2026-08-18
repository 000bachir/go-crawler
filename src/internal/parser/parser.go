package parser

import (
	"fmt"
	"io"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
)

type GoqueryParser struct {
	excludeExtension map[string]bool
}

func NewGoqueryParser() *GoqueryParser {
	return &GoqueryParser{
		excludeExtension: make(map[string]bool),
	}
}

// this function will silently drop and skip extensions like .pdf , .zip so on and so forht

func (p *GoqueryParser) ExcludeExtension(exts ...string) {
	for _, extension := range exts {
		p.excludeExtension[extension] = true
	}
}

// parse reads html from the reader , resolve all links based on the baseurl
// return a list of valid urls that are absolute found in the page
func (p *GoqueryParser) Parse(baseUrl string, reader io.Reader) ([]*url.URL, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("parsing html : %w", err)
	}
	return p.ExtractLinks(doc, baseUrl), nil
}

func (p *GoqueryParser) ExtractLinks(doc *goquery.Document, baseURL string) []*url.URL {
	var links []*url.URL

	elements := doc.Find("a , link")
	for i := 0; i < elements.Length(); i++ {
		element := elements.Eq(i)

		href, ok := p.HrefFromElement(element)
		if !ok {
			continue
		}
		link, ok := resolveRelativeURL(baseURL, href)
		if !ok {
			continue
		}

		links = append(links, link)
	}
	return links
}

func (p *GoqueryParser) HrefFromElement(s *goquery.Selection) (string, bool) {
	switch goquery.NodeName(s) {
	case "a":
		return p.HrefFromAnchor(s)
	case "link":
		return p.HrefFromLinkTag(s)
	default:
		return "", false

	}
}

func (p *GoqueryParser) HrefFromAnchor(s *goquery.Selection) (string, bool) {
	href, exists := s.Attr("href")
	if !exists {
		return "", false
	}
	if p.isExcludeExtension(href) {
		return "", false
	}
	return href, true
}

func (p *GoqueryParser) HrefFromLinkTag(s *goquery.Selection) (string, bool) {
	rel, exists := s.Attr("ref")

	if !exists || rel != "canonical" {
		return "", false
	}

	href, exists := s.Attr("href")
	if !exists {
		return "", false
	}

	if p.isExcludeExtension(href) {
		return "", false
	}
	return href, true
}

func (p *GoqueryParser) isExcludeExtension(href string) bool {
	u, err := url.Parse(href)
	if err != nil {
		return false
	}
	return p.excludeExtension[path.Ext(u.Path)]
}
func resolveRelativeURL(baseRaw, refRaw string) (*url.URL, bool) {
	base, err := url.Parse(baseRaw)
	if err != nil {
		return nil, false
	}

	ref, err := url.Parse(refRaw)
	if err != nil {
		return nil, false
	}

	resolved := base.ResolveReference(ref)

	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return nil, false
	}

	// Strip the fragment: /page#section1 and /page#section2 are the same resource.
	resolved.Fragment = ""

	return resolved, true
}
