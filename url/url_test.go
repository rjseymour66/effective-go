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
	const in = "foo.com:80"

	u := &URL{Host: in}
	if got, want := u.Port(), "80"; got != want {
		t.Errorf("for host %q; got %q, want %q", in, got, want)
	}
}
