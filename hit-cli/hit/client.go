package hit

import (
	"context"
	"net/http"
	"time"
)

// Client sends HTTP requests and returns an aggregated performance
// result. The fields should not be changed after initializing.
type Client struct {
	C   int // C is the concurrency level
	RPS int // RPS throttles the requests per second
}

// Do sends an HTTP request and returns an aggregated result.
func (c *Client) Do(r *http.Request, n int) *Result {
	t := time.Now()
	sum := c.do(r, n)
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(r *http.Request, n int) *Result {
	p := produce(n, func() *http.Request {
		return r.Clone(context.TODO())
	})
	if c.RPS > 0 {
		p = throttle(p, time.Second/time.Duration(c.RPS*c.C))
	}
	var sum Result
	for result := range split(p, c.C, Send) {
		sum.Merge(result)
	}
	return &sum
}
