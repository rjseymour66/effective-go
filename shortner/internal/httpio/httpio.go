package httpio

import (
	"encoding/json"
	"io"
	"net/http"
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
