package user

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
	r.Mux.HandleFunc("GET /user", r.Controller.FindAll)
	r.Mux.HandleFunc("GET /user/{username}", r.Controller.FindByUsername)
	r.Mux.HandleFunc("POST /user", r.Controller.Create)
	r.Mux.HandleFunc("PATCH /user/{id}", r.Controller.Update)
	r.Mux.HandleFunc("DELETE /user/{id}", r.Controller.Delete)
}
