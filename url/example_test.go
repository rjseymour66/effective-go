package url_test

import (
	"fmt"
	"log"

	"github.com/rjseymour66/effective-go/url"
)

func ExampleURL() {
	u, err := url.Parse("http://foo.com/go")
	if err != nil {
		log.Fatal(err)
	}

	u.Scheme = "https"
	fmt.Println(u)
	// Output: https://foo.com/go
}
