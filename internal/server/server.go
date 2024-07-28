package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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
	splittedCorsListEnv := strings.Split(os.Getenv("ALLOWED_CORS"), ",")
	allowedCorsList := make([]string, 0)
	allowedCorsList = append(allowedCorsList, splittedCorsListEnv...)
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedCorsList,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	err := http.ListenAndServe(fmt.Sprintf(":%v", s.Port), c.Handler(s.Mux))
	return err
}
