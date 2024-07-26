package server

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

type Server struct {
	Mux  *http.ServeMux
	Port int
}

func New(mux *http.ServeMux, port int) *Server {
	return &Server{
		Mux:  mux,
		Port: port,
	}
}

func (s *Server) Run() error {
	err := http.ListenAndServe(fmt.Sprintf(":%v", s.Port), cors.Default().Handler(s.Mux))
	return err
}
