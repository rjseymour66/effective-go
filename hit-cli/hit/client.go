package hit

import (
	"context"
	"net/http"
	"time"
)

// Client sends HTTP requests and returns an aggregated performance
// result. The fields should not be changed after initializing.
type Client struct {
	C       int // C is the concurrency level
	RPS     int // RPS throttles the requests per second
	Timeout time.Duration
}

// Do sends an HTTP request and returns an aggregated result.
func (c *Client) Do(ctx context.Context, r *http.Request, n int) *Result {
	t := time.Now()
	sum := c.do(ctx, r, n)
	return sum.Finalize(time.Since(t))
}

// clone the request to let the producer produce identical request values to
// send HTTP requests to the same URL. If request values were not cloned, there
// would be clashes with other req values bc each req is stateful.
// https://pkg.go.dev/net/http#Request.Clone
func (c *Client) do(ctx context.Context, r *http.Request, n int) *Result {
	p := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})
	if c.RPS > 0 {
		p = throttle(p, time.Second/time.Duration(c.RPS*c.C))
	}
	var (
		sum    Result
		client = c.client()
	)
	defer client.CloseIdleConnections()
	for result := range split(p, c.C, c.send(client)) {
		sum.Merge(result)
	}
	return &sum
}

func (c *Client) send(client *http.Client) SendFunc {
	return func(r *http.Request) *Result {
		return Send(client, r)
	}
}

func (c *Client) client() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: c.C,
		},
	}
}
