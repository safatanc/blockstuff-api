package minecraftserver

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
	r.Mux.HandleFunc("GET /minecraftserver", r.Controller.FindAll)
	r.Mux.HandleFunc("GET /minecraftserver/{slug}", r.Controller.FindBySlug)
	r.Mux.Handle("GET /minecraftserver/{slug}/detail", r.Middleware.Auth(http.HandlerFunc(r.Controller.FindBySlugDetail)))
	r.Mux.Handle("POST /minecraftserver", r.Middleware.Auth(http.HandlerFunc(r.Controller.Create)))
	r.Mux.Handle("PATCH /minecraftserver/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Update)))
	r.Mux.Handle("PATCH /minecraftserver/{id}/rcon", r.Middleware.Auth(http.HandlerFunc(r.Controller.UpdateRcon)))
	r.Mux.Handle("DELETE /minecraftserver/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Delete)))
}
