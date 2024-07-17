package minecraftserver

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
	r.Mux.HandleFunc("GET /minecraftserver", r.Controller.FindAll)
	r.Mux.HandleFunc("GET /minecraftserver/{ip}", r.Controller.FindByIP)
	r.Mux.HandleFunc("POST /minecraftserver", r.Controller.Create)
	r.Mux.HandleFunc("PATCH /minecraftserver/{id}", r.Controller.Update)
	r.Mux.HandleFunc("DELETE /minecraftserver/{id}", r.Controller.Delete)
}
