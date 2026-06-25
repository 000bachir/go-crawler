package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func newServer(handler http.HandlerFunc) *httptest.Server {
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
	test := newServer(htmlHandler(simpleHTML))
	defer test.Close()

	fetch := NewHTTPFetcher()
	doc, err := fetch.Fetch(context.Background(), test.URL)
	if err != nil {
		t.Fatal("expected no error got :", err)
	}
	if doc == nil {
		t.Fatal("expected parsed document got nil ", doc)
	}
	if doc.Type != html.DocumentNode {
		t.Fatal("expected document node at root got: ", doc.Type)
	}
}

func TestFetch_WrongContentType(t *testing.T) {
	srv := newServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key":"value"}`))
	})
	defer srv.Close()

	f := NewHTTPFetcher()
	_, err := f.Fetch(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected an error for non-HTML content type, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected format") {
		t.Errorf("error message should mention unexpected format, got: %v", err)
	}
}

func TestFetchCancelledContext(t *testing.T) {
	srv := newServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key" :"value" }`))
	})
	defer srv.Close()
	fetcher := NewHTTPFetcher()
	_, err := fetcher.Fetch(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expecpted error for non html content type , got nil")
	}
	if !strings.Contains(err.Error(), "unexpected format") {
		t.Errorf("expected 'unexpected format' in error , got : %v", err)
	}
}
