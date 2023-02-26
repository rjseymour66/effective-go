package hit

import (
	"net/http"
	"time"
)

// Client sends HTTP requests and returns an aggregated performance
// result. The fields should not be changed after initializing.
type Client struct {
	// TODO
}

// Do sends an HTTP request and returns an aggregated result.
func (c *Client) Do(r *http.Request, n int) *Result {
	t := time.Now()
	sum := c.do(r, n)
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(r *http.Request, n int) *Result {
	var sum Result
	for ; n > 0; n-- {
		sum.Merge(Send(r))
	}
	return &sum
}
