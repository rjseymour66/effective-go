package short

import (
	"effective-go/shortner/internal/httpio"
	"effective-go/shortner/linkit"
	"fmt"
	"net/http"
)

const (
	shorteningRoute  = "/s"
	resolveRoute     = "/r/"
	healthCheckRoute = "/health"
)

// mux is an unexported http.Handler
type mux http.Handler

// Server is a custom server type.
type Server struct {
	mux // the server only exports ServeHTTP
}

// NewServer returns an instance of the custom Server type.
func NewServer() *Server {
	var s Server
	s.registerRoutes()
	return &s
}

func (s *Server) registerRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc(shorteningRoute, s.shorteningHandler)
	mux.HandleFunc(resolveRoute, s.resolveHandler)
	mux.HandleFunc(healthCheckRoute, s.healthCheckHandler)
	s.mux = mux
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

// the shortening handler decodes the client's JSON request by reading the
// Request.Body (an io.Reader) and storing it in the input variable. After
// processing the request, it responds with the shortened key in JSON format.
func (s *Server) shorteningHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var input struct {
		URL string
		Key string
	}
	err := httpio.Decode(http.MaxBytesReader(w, r.Body, 4_096), &input)
	if err != nil {
		http.Error(w, "cannot decode JSON", http.StatusBadRequest)
		return
	}

	ln := link{
		uri:      input.URL,
		shortKey: input.Key,
	}

	if err := checkLink(ln); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = httpio.Encode(w, http.StatusCreated, map[string]any{
		"key": ln.shortKey,
	})
}

func (s *Server) resolveHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(resolveRoute):]

	if err := checkShortKey(key); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// use dummy data for now and carelessly expose internal details
	if key == "fortesting" {
		http.Error(w, "db at IP ... failed", http.StatusInternalServerError)
		return
	}
	if key != "go" {
		http.Error(w, linkit.ErrNotExist.Error(), http.StatusNotFound)
		return
	}
	const uri = "https://go.dev"
	http.Redirect(w, r, uri, http.StatusFound)
}
