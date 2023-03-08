package hit

import (
	"net/http"
	"time"
)

// Produce calls fn n times and sends results to out.
func Produce(out chan<- *http.Request, n int, fn func() *http.Request) {
	for ; n > 0; n-- {
		out <- fn()
	}
}

// produce runs Produce in a goroutine.
func produce(n int, fn func() *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Produce(out, n, fn)
	}()
	return out
}

// Throttle slows down receiving from in by delay and
// sends what it receives from in to out.
func Throttle(in <-chan *http.Request, out chan<- *http.Request, delay time.Duration) {
	t := time.NewTicker(delay)
	defer t.Stop()

	for r := range in {
		<-t.C
		out <- r
	}
}

// throttle runs Throttle in a goroutine.
func throttle(in <-chan *http.Request, delay time.Duration) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Throttle(in, out, delay)
	}()
	return out
}
