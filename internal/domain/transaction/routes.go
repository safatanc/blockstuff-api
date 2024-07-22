package transaction

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
	r.Mux.Handle("GET /transaction", r.Middleware.Auth(http.HandlerFunc(r.Controller.FindAll)))
	r.Mux.HandleFunc("GET /transaction/{code}", r.Controller.FindByCode)
	r.Mux.Handle("POST /transaction/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Create)))
	r.Mux.Handle("PATCH /transaction/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Update)))
	r.Mux.Handle("DELETE /transaction/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Delete)))
}
