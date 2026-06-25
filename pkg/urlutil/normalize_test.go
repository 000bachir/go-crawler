package urlutil

import (
	"net/url"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		base     string
		want     string
		wantOK   bool
	}{
		{
			name:     "relative path resolved against base",
			inputURL: "/foo",
			base:     "https://example.com/bar",
			want:     "https://example.com/foo",
			wantOK:   true,
		},
		{
			name:     "absolute https URL",
			inputURL: "https://other.com/page",
			base:     "https://example.com",
			want:     "https://other.com/page",
			wantOK:   true,
		},
		{
			name:     "absolute http URL",
			inputURL: "http://other.com/page",
			base:     "https://example.com",
			want:     "http://other.com/page",
			wantOK:   true,
		},
		{
			name:     "fragment removed",
			inputURL: "/foo#section",
			base:     "https://example.com",
			want:     "https://example.com/foo",
			wantOK:   true,
		},
		{
			name:     "unsupported scheme",
			inputURL: "ftp://example.com/file",
			base:     "https://example.com",
			want:     "",
			wantOK:   false,
		},
		{
			name:     "javascript scheme",
			inputURL: "javascript:alert(1)",
			base:     "https://example.com",
			want:     "",
			wantOK:   false,
		},
		{
			name:     "invalid base URL",
			inputURL: "/foo",
			base:     "://invalid",
			want:     "",
			wantOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := url.Parse(tt.inputURL)
			if err != nil {
				t.Errorf("failed to parse the given url : %v", err)
			}

			got, ok := Normalize(url, tt.base)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v , want %v", ok, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("got %q , want %q", got, tt.want)
			}

		})
	}

}
