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
	mux.Handle(shorteningRoute, httpio.Handler(s.shorteningHandler))
	mux.Handle(resolveRoute, httpio.Handler(s.resolveHandler))
	mux.HandleFunc(healthCheckRoute, s.healthCheckHandler)
	s.mux = mux
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

// the shortening handler decodes the client's JSON request by reading the
// Request.Body (an io.Reader) and storing it in the input variable. After
// processing the request, it responds with the shortened key in JSON format.
func (s *Server) shorteningHandler(w http.ResponseWriter, r *http.Request) http.Handler {
	if r.Method != http.MethodPost {
		return httpio.Error(http.StatusMethodNotAllowed, "method not allowed")
	}
	var input struct {
		URL string
		Key string
	}
	err := httpio.Decode(http.MaxBytesReader(w, r.Body, 4_096), &input)
	if err != nil {
		return httpio.Error(http.StatusBadRequest, "cannot decode JSON")
	}

	ln := link{
		uri:      input.URL,
		shortKey: input.Key,
	}

	if err := checkLink(ln); err != nil {
		return httpio.Error(http.StatusBadRequest, err.Error())
	}

	_ = httpio.Encode(w, http.StatusCreated, map[string]any{
		"key": ln.shortKey,
	})
	return httpio.JSON(http.StatusCreated, map[string]any{
		"key": ln.shortKey,
	})
}

func (s *Server) resolveHandler(w http.ResponseWriter, r *http.Request) http.Handler {
	key := r.URL.Path[len(resolveRoute):]

	if err := checkShortKey(key); err != nil {
		return httpio.Error(http.StatusBadRequest, err.Error())
	}
	// use dummy data for now and carelessly expose internal details
	if key == "fortesting" {
		return httpio.Error(http.StatusInternalServerError, "db at IP ... failed")
	}
	if key != "go" {
		return httpio.Error(http.StatusNotFound, linkit.ErrNotExist.Error())
	}
	const uri = "https://go.dev"
	http.Redirect(w, r, uri, http.StatusFound)

	return nil
}
