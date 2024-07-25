package item

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
	r.Mux.HandleFunc("GET /minecraftserver/{minecraft_server_id}/item", r.Controller.FindAll)
	r.Mux.HandleFunc("GET /minecraftserver/{minecraft_server_id}/item/{slug}", r.Controller.FindBySlug)
	r.Mux.Handle("POST /minecraftserver/{minecraft_server_id}/item", r.Middleware.Auth(http.HandlerFunc(r.Controller.Create)))
	r.Mux.Handle("POST /minecraftserver/{minecraft_server_id}/item/{id}/image", r.Middleware.Auth(http.HandlerFunc(r.Controller.AddImage)))
	r.Mux.Handle("POST /minecraftserver/{minecraft_server_id}/item/{id}/action", r.Middleware.Auth(http.HandlerFunc(r.Controller.AddAction)))
	r.Mux.Handle("PATCH /minecraftserver/{minecraft_server_id}/item/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Update)))
	r.Mux.Handle("DELETE /minecraftserver/{minecraft_server_id}/item/{id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.Delete)))
	r.Mux.Handle("DELETE /minecraftserver/{minecraft_server_id}/item/{id}/image/{item_image_id}", r.Middleware.Auth(http.HandlerFunc(r.Controller.DeleteImage)))
}
