package httpio

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Decode reads from a reader into any value
func Decode(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// Encode writes a Go value as JSON to the client
func Encode(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

// LoggingMiddleware logs each handler response time.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Since(start)
		Log(r.Context(), "%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, end)
	})
}

// Log is the custom logger implementation.
func Log(ctx context.Context, format string, args ...any) {
	s, _ := ctx.Value(http.ServerContextKey).(*http.Server)
	if s == nil || s.ErrorLog == nil {
		return
	}
	s.ErrorLog.Printf(format, args...)
}
