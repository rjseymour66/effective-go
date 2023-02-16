package url

import (
	"errors"
	"strings"
)

// A URL represents a parsed URL.
type URL struct {
	// https://foo.com
	Scheme string // https
	Host   string // foo
	Path   string // go
}

// Parse takes a raw URL string and returns a URL struct.
func Parse(rawurl string) (*URL, error) {
	i := strings.Index(rawurl, "://")
	if i < 0 {
		return nil, errors.New("missing scheme")
	}

	scheme, rest := rawurl[:i], rawurl[i+3:]
	host, path := rest, ""
	if i := strings.Index(rest, "/"); i >= 0 {
		host, path = rest[:i], rest[i+1:]
	}

	return &URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}, nil
}

// Port returns the port number from the URL Host field.
func (u *URL) Port() string {
	return ""
}
