package main

import (
	"effective-go/shortner/short"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {

	const (
		addr    = "localhost:8080"
		timeout = 10 * time.Second
	)

	fmt.Fprintln(os.Stderr, "starting the server on", addr)

	// shortener := http.HandlerFunc(
	// 	func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintln(w, "hello from the shortener server!")
	// 		// w.Write([]byte("hello from the shortener server!"))
	// 	},
	// )

	shortener := short.NewServer()

	server := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(shortener, timeout, "timeout"),
		ReadTimeout: timeout,
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintln(os.Stderr, "server closed unexpectedly:", err)
	}
}
