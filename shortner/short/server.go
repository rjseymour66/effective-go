package short

import (
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
	// http.Handler
	mux // the server only exports ServeHTTP
}

// type Server struct {
// 	mux *http.ServeMux
// }

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

// func (s *Server) registerRoutes() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc(shorteningRoute, s.shorteningHandler)
// 	mux.HandleFunc(resolveRoute, s.resolveHandler)
// 	mux.HandleFunc(healthCheckRoute, s.healthCheckHandler)
// 	s.mux = mux
// }

// func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	s.mux.ServeHTTP(w, r)
// }

// func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch p := r.URL.Path; {
// 	case p == healthCheckRoute:
// 		s.healthCheckHandler(w, r)
// 	case strings.HasPrefix(p, resolveRoute):
// 		s.resolveHandler(w, r)
// 	case strings.HasPrefix(p, shorteningRoute):
// 		s.shorteningHandler(w, r)
// 	default:
// 		http.NotFound(w, r)
// 	}
// }

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func (s *Server) shorteningHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "go")
}

func (s *Server) resolveHandler(w http.ResponseWriter, r *http.Request) {
	const uri = "https://go.dev"
	http.Redirect(w, r, uri, http.StatusFound)
}
