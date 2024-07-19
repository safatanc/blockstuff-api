package user

import (
	"net/http"

	"github.com/safatanc/blockstuff-api/internal/middleware"
)

type Routes struct {
	Mux        *http.ServeMux
	Controller *Controller
	Middleware *middleware.Middleware
}

func NewRoutes(mux *http.ServeMux, controller *Controller, middleware *middleware.Middleware) *Routes {
	return &Routes{
		Mux:        mux,
		Controller: controller,
		Middleware: middleware,
	}
}

func (r *Routes) Init() {
	r.Mux.Handle("GET /user", r.Middleware.Auth(http.HandlerFunc(r.Controller.FindAll)))
	r.Mux.Handle("GET /user/{username}", r.Middleware.Auth(http.HandlerFunc(r.Controller.FindByUsername)))
	r.Mux.Handle("POST /user", r.Middleware.Auth(http.HandlerFunc(r.Controller.Create)))
	r.Mux.Handle("PATCH /user/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Update)))
	r.Mux.Handle("DELETE /user/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Delete)))
}
