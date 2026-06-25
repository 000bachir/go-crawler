package main

// what it should do : given a url return the raw data
import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Attributes struct {
	Key, Value string
}

type NodeInfo struct {
	Tag        string
	Attributes []Attributes
	Depth      int
	ParentTag  string
	URL        string
	Text       string
}

type HTTPFetcher struct {
	client *http.Client
}

func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, redirects []*http.Request) error {
				if len(redirects) >= 5 {
					return fmt.Errorf("to manny redirects (>= %d) for url : %s", len(redirects), req.URL)
				}
				return nil
			},
		},
	}
}

func (f *HTTPFetcher) Fetch(ctx context.Context, url string) (*html.Node, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request %w", err)
	}
	// setting up the header
	request.Header.Set("User-Agent", "crawlox")

	response, err := f.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making the request %w", err)
	}
	// checking the type of the content
	content := response.Header.Get("Content-Type")

	if !strings.Contains(content, "text/html") {
		response.Body.Close()
		return nil, fmt.Errorf("unexpected format returned got : %q", content)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("warning : error closing the  request %v", err)
		}
	}()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}
	const MaxBodyBytes int = 5 * 1024 * 1024

	limit := io.LimitReader(response.Body, int64(MaxBodyBytes))

	doc, err := html.Parse(limit)
	if err != nil {
		return nil, fmt.Errorf("parsing the html %w", err)
	}

	return doc, nil
}

func (f *HTTPFetcher) FetchWithRetry(ctx context.Context, url string, maxRetry int) (*html.Node, error) {
	var err error
	for i := range maxRetry {
		var doc *html.Node
		doc, err = f.Fetch(ctx, url)
		if err == nil {
			return doc, nil
		}

		wait := time.Duration(1<<i) * time.Second
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-time.After(wait):
		}
	}
	return nil, fmt.Errorf("all %d retries failed: %w", maxRetry, err)
}

func ExtractAttributes(node *html.Node) []Attributes {
	attributes := make([]Attributes, 0, len(node.Attr))
	for _, attr := range node.Attr {
		attributes = append(attributes, Attributes{
			Key:   attr.Key,
			Value: attr.Val,
		})
	}
	return attributes
}

func TraversHTMLResponse(root *html.Node, maxDepth int) []NodeInfo {
	type stackEnty struct {
		node  *html.Node
		depth int
	}
	var resutls []NodeInfo
	stack := []stackEnty{
		{
			node:  root,
			depth: 0,
		},
	}
	for len(stack) > 0 {
		entry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		node, depth := entry.node, entry.depth
		if maxDepth > 0 && depth > maxDepth {
			continue
		}

		if node.Type == html.ElementNode {
			resutls = append(resutls, NodeInfo{
				Tag:        node.Data,
				Attributes: ExtractAttributes(node),
				Depth:      depth,
				URL:        "",
			})
		}
		for child := node.LastChild; child != nil; child = child.PrevSibling {
			stack = append(stack, stackEnty{
				node: child, depth: depth + 1,
			})
		}
	}
	return resutls
}
