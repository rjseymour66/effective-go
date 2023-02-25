package main

import "fmt"

const (
	bannerText = `
 __  __    ___    _______
/\ \_\ \  /\  \  /\__   __\
\ \  __ \ \ \  \ \/_/\  \/
 \ \_\ \_\ \ \__\   \ \__\
  \/ /\/_/  \/__/    \/__/
`

	usageText = `
Usage:
	-url
		HTTP server URL to make requests (required)
	-n
		Number of requests to make
	-c
		Concurrency level`
)

func banner() string { return bannerText[1:] }
func usage() string  { return usageText[1:] }

func main() {
	fmt.Println(banner())
	fmt.Println(usage())
}
