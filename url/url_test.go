package url

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	const rawurl = "https://foo.com/go"

	u, err := Parse(rawurl)
	if err != nil {
		t.Fatalf("Parse(%q) err = %q, want nil", rawurl, err)
	}

	want := "https"
	got := u.Scheme
	if got != want {
		t.Errorf("Parse(%q).Scheme = %q, want %q", rawurl, got, want)
	}

	if got, want := u.Host, "foo.com"; got != want {
		t.Errorf("Parse(%q).Host = %q, want %q", rawurl, got, want)
	}

	if got, want := u.Path, "go"; got != want {
		t.Errorf("Parse(%q).Path = %q, want %q", rawurl, got, want)
	}
}

func TestParseInvalidURLs(t *testing.T) {
	tests := map[string]string{
		"missing scheme": "foo.com",
		"empty scheme":   "://foo.com",
	}
	for name, in := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse(in); err == nil {
				t.Errorf("Parse(%q)=nil, want an error", in)
			}
		})
	}
}

func TestURLHost(t *testing.T) {
	tests := map[string]struct {
		in       string
		hostname string
		port     string
	}{
		"with port":       {in: "foo.com:80", hostname: "foo.com", port: "80"},
		"empty port":      {in: "foo.com:", hostname: "foo.com", port: ""},
		"without port":    {in: "foo.com", hostname: "foo.com", port: ""},
		"ip with port":    {in: "1.2.3.4:90", hostname: "1.2.3.4", port: "90"},
		"ip without port": {in: "1.2.3.4", hostname: "1.2.3.4", port: ""},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("Port/%s/%s", name, tt.in), func(t *testing.T) {
			u := &URL{Host: tt.in}
			if got, want := u.Port(), tt.port; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
		t.Run(fmt.Sprintf("Hostname/%s/%s", name, tt.in), func(t *testing.T) {
			u := &URL{Host: tt.in}
			if got, want := u.Hostname(), tt.hostname; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

func TestURLString(t *testing.T) {
	tests := map[string]struct {
		url  *URL
		want string
	}{
		"nil url":   {url: nil, want: ""},
		"empty url": {url: &URL{}, want: ""},
		"sheme":     {url: &URL{Scheme: "https"}, want: "https://"},
		"host": {
			url:  &URL{Scheme: "https", Host: "foo.com"},
			want: "https://foo.com",
		},
		"path": {
			url:  &URL{Scheme: "https", Host: "foo.com", Path: "go"},
			want: "https://foo.com/go",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if g, w := tt.url, tt.want; g.String() != w {
				t.Errorf("url: %#v\ngot:  %q\nwant: %q", g, g, w)
			}
		})
	}
}

func BenchmarkURLString(b *testing.B) {
	b.ReportAllocs()
	b.Logf("Loop %d times\n", b.N)

	u := &URL{Scheme: "https", Host: "foo.com", Path: "go"}

	for i := 0; i < b.N; i++ {
		u.String()
	}
}

// func BenchmarkURLString(b *testing.B) {
// 	var benchmarks = []*URL{
// 		{Scheme: "https"},
// 		{Scheme: "https", Host: "foo.com"},
// 		{Scheme: "https", Host: "foo.com", Path: "go"},
// 	}
// 	for _, u := range benchmarks {
// 		b.Run(u.String(), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				u.String()
// 			}
// 		})
// 	}
// }
