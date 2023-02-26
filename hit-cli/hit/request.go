package hit

import (
	"fmt"
	"net/http"
	"time"
)

// SendFunc is the type of the function called by Client.Do
// to send an HTTP request and return a performance result.
type SendFunc func(*http.Request) *Result

// Send an HTTP request and return a performance result.
func Send(r *http.Request) *Result {
	t := time.Now()

	fmt.Printf("request: %s\n", r.URL)
	time.Sleep(100 * time.Millisecond)

	return &Result{
		Duration: time.Since(t),
		Bytes:    10,
		Status:   http.StatusOK,
	}
}
