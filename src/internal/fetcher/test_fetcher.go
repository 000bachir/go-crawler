package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/html"
)

func newServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func htmlHandler(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}

const simpleHTML = `<!DOCTYPE html>
<html>
  <head><title>Test</title></head>
  <body>
    <h1 class="title">Hello</h1>
    <a href="https://example.com">Link</a>
  </body>
</html>`

// fetch

func TestFetchSuccess(t *testing.T) {
	srv := newServer(htmlHandler(simpleHTML))
	defer srv.Close()

	fetcher := NewHTTPFetcher()
	doc, err := fetcher.Fetch(context.Background(), srv.URL)
	if err != nil {
		t.Fatal("expected no error got %v", err)
	}
	if doc == nil {
		t.Fatal("expected a parsed document got nil")
	}
	if doc.Type != html.DocumentNode {
		t.Errorf("expected document node at root , got %v", doc.Type)
	}
}
