package main

import (
	"effective-go/hit-cli/hit"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

const bannerText = `
 __  __    ___    ________
/\ \_\ \  /\  \  /\__   __\
\ \  __ \ \ \  \ \/_/\  \_/
 \ \_\ \_\ \ \__\   \ \__\
  \/_/\/_/  \/__/    \/__/
`

func banner() string { return bannerText[1:] }

func main() {
	if err := run(flag.CommandLine, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error occurred:", err)
		os.Exit(1)
	}
}

func run(s *flag.FlagSet, args []string, out io.Writer) error {
	f := &flags{
		n: 100,
		c: runtime.NumCPU(),
	}
	if err := f.parse(s, args); err != nil {
		return err
	}

	fmt.Fprintln(out, banner())
	fmt.Fprintf(out, "Making %d requests to %s with a concurrency level of %d.\n",
		f.n, f.url, f.c)

	request, err := http.NewRequest(http.MethodGet, f.url, http.NoBody)
	if err != nil {
		return err
	}
	var c hit.Client
	sum := c.Do(request, f.n)
	sum.Fprint(out)

	return nil
}
