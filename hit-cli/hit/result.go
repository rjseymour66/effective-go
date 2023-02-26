package hit

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Result is a request's result.
type Result struct {
	RPS      float64
	Requests int
	Errors   int
	Bytes    int64
	Duration time.Duration
	Fastest  time.Duration
	Slowest  time.Duration
	Status   int
	Error    error
}

// Merge this result with another
func (r *Result) Merge(o *Result) {
	r.Requests++
	r.Bytes += o.Bytes

	if r.Fastest == 0 || o.Duration < r.Fastest {
		r.Fastest = o.Duration
	}
	if o.Duration > r.Slowest {
		r.Slowest = o.Duration
	}

	switch {
	case o.Error != nil:
		fallthrough
	case o.Status >= http.StatusBadRequest:
		r.Errors++
	}
}

// Finalize the total duration and calculate RPS.
func (r *Result) Finalize(total time.Duration) *Result {
	r.Duration = total
	r.RPS = float64(r.Requests) / total.Seconds()
	return r
}

// Fprint the result to an io.Writer
func (r *Result) Fprint(out io.Writer) {
	p := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}
	p("\nSummary:\n")
	p("\tSuccess		: %0.f%%\n", r.success())
	p("\tRPS		: %.1f\n", r.RPS)
	p("\tRequests	: %d\n", r.Requests)
	p("\tErrors		: %d\n", r.Errors)
	p("\tBytes		: %d\n", r.Bytes)
	p("\tDuration	: %s\n", round(r.Duration))
	if r.Requests > 1 {
		p("\tFastest		: %s\n", round(r.Fastest))
		p("\tSlowest		: %s\n", round(r.Slowest))
	}
}

func (r *Result) success() float64 {
	rr, e := float64(r.Requests), float64(r.Errors)
	return (rr - e) / rr * 100
}

func round(t time.Duration) time.Duration {
	return t.Round(time.Microsecond)
}

// String implements the Stringer interface
// func (r *Result) String() string {
// 	var s strings.Builder
// 	r.Fprint(&s)
// 	return s.String()
// }
