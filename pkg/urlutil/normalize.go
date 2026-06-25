package urlutil

import "net/url"

func Normalize(inputUrl *url.URL, raw string) (string, bool) {
	base, err := url.Parse(raw)
	if err != nil {
		return "", false
	}
	resolved := base.ResolveReference(inputUrl)
	switch resolved.Scheme {
	case "http", "https":
	default:
		return "", false
	}
	resolved.Fragment = ""
	return resolved.String(), true
}
