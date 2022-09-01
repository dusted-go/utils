package typ

import "strings"

type URL string

// NewURL creates a new URL object.
func NewURL(v string) URL {
	if len(v) == 0 {
		return URL("")
	}
	if !(strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://")) {
		v = "https://" + v
	}
	return URL(v)
}

// Empty checks if the URL is empty.
func (u URL) Empty() bool {
	return len(string(u)) > 0
}

// String returns the URL as raw string.
func (u URL) String() string {
	return string(u)
}

// Pretty returns a string without any http/https prefix.
func (u URL) Pretty() string {
	return strings.TrimPrefix(
		strings.TrimPrefix(string(u),
			"https://"),
		"http://")
}
