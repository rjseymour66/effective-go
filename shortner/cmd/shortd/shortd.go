package main

import (
	"effective-go/shortner/internal/httpio"
	"effective-go/shortner/short"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	const (
		addr    = "localhost:8080"
		timeout = 10 * time.Second
	)

	logger := log.New(os.Stderr, "shortner: ", log.LstdFlags|log.Lmsgprefix)
	logger.Println("starting the server on", addr)

	fmt.Fprintln(os.Stderr, "starting the server on", addr)

	shortener := short.NewServer()

	server := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(shortener, timeout, "timeout"),
		ReadTimeout: timeout,
	}

	if os.Getenv("LINKIT_DEBUG") == "1" {
		server.ErrorLog = logger
		server.Handler = httpio.LoggingMiddleware(server.Handler)
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintln(os.Stderr, "server closed unexpectedly:", err)
	}
}
