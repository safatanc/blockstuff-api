package auth

import "net/http"

type Routes struct {
	Mux        *http.ServeMux
	Controller *Controller
}

func NewRoutes(mux *http.ServeMux, controller *Controller) *Routes {
	return &Routes{
		Mux:        mux,
		Controller: controller,
	}
}

func (r *Routes) Init() {
	r.Mux.HandleFunc("POST /auth/login", r.Controller.Login)
	r.Mux.HandleFunc("POST /auth/register", r.Controller.Register)
}
