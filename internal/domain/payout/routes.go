package payout

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
	r.Mux.Handle("GET /payout", r.Middleware.Auth(http.HandlerFunc(r.Controller.FindAll)))
	r.Mux.Handle("GET /payout/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.FindByID)))
	r.Mux.Handle("POST /payout", r.Middleware.Auth(http.HandlerFunc(r.Controller.Create)))
	r.Mux.Handle("PATCH /payout/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Update)))
	r.Mux.Handle("DELETE /payout/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Delete)))
	r.Mux.Handle("GET /payout/channel", r.Middleware.Auth(http.HandlerFunc(r.Controller.FindPayoutChannels)))
	r.Mux.Handle("GET /payout/channel/user/{username}", r.Middleware.Auth(http.HandlerFunc(r.Controller.GetPayoutChannel)))
	r.Mux.Handle("PUT /payout/channel/user/{username}", r.Middleware.Auth(http.HandlerFunc(r.Controller.SetPayoutChannel)))
}
