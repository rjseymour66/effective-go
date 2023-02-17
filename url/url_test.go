package url

import "testing"

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

func TestURLPort(t *testing.T) {
	tests := map[string]struct {
		in   string
		port string
	}{
		"with port":       {in: "foo.com:80", port: "80"},
		"with empty port": {in: "foo.com:", port: ""},
		"without port":    {in: "foo.com", port: ""},
		"ip with port":    {in: "1.2.3.4:90", port: "90"},
		"ip without port": {in: "1.2.3.4", port: ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			u := &URL{Host: tt.in}
			if got, want := u.Port(), tt.port; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}
