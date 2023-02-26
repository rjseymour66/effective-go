package main

import (
	"flag"
	"fmt"
	"io"
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

	fmt.Println(banner())
	fmt.Printf("Making %d requests to %s with a concurrency level of %d.\n",
		f.n, f.url, f.c)

	// hit pkg integration here

	return nil
}
