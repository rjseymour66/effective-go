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
	if i < 1 {
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

// Port returns the port number from u.HOst, without the leading colon.
func (u *URL) Port() string {
	i := strings.Index(u.Host, ":")
	if i < 0 {
		return ""
	}
	return u.Host[i+1:]
}

// Hostname returns u.Host, and strips any port number, if present.
func (u *URL) Hostname() string {
	i := strings.Index(u.Host, ":")
	if i < 0 {
		return u.Host
	}
	return u.Host[:i]
}

// String reassembles the URL into a URL string.
func (u *URL) String() string {
	if u == nil {
		return ""
	}
	var s strings.Builder
	if sc := u.Scheme; sc != "" {
		s.WriteString(sc)
		s.WriteString("://")
	}
	if h := u.Host; h != "" {
		s.WriteString(h)
	}
	if p := u.Path; p != "" {
		s.WriteByte('/')
		s.WriteString(p)
	}
	return s.String()
}
