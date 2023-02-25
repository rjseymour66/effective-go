package main

import (
	"fmt"
	"os"
	"runtime"
)

const bannerText = `
 __  __    ___    _______
/\ \_\ \  /\  \  /\__   __\
\ \  __ \ \ \  \ \/_/\  \/
 \ \_\ \_\ \ \__\   \ \__\
  \/ /\/_/  \/__/    \/__/
`

func banner() string { return bannerText[1:] }

func main() {
	f := &flags{
		n: 100,
		c: runtime.NumCPU(),
	}
	if err := f.parse(); err != nil {
		os.Exit(1)
	}

	fmt.Println(banner())
	fmt.Printf("Making %d requests to %s with a concurrency level of %d.\n",
		f.n, f.url, f.c)
}
