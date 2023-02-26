package hit

import "net/http"

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
