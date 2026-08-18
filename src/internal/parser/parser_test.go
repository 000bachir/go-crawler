package parser_test

import (
	"errors"
	"net/url"
	"testing"
)

const testBaseURl = "https://exemple.com"

func containsURL(links []*url.URL, target string) bool {
	for _, l := range links {
		if l.String() == target {
			return true
		}
	}
	return false
}

type ErrorReader struct{}

func (ErrorReader) Read(_ []byte) (int, error) {
	return 0, errors.New("simulated read error")
}

func TestResolveRelativeURL_AbsoluteURL(t *testing.T) {
	got, ok := resolveRelativeURL(testBaseURl, "https://other.com/page")
	if !ok {
		t.Fatal("expected ok = true for an absolute url")
	}
}
