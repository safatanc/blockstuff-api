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
	r.Mux.HandleFunc("POST /auth/verify", r.Controller.Verify)
	r.Mux.HandleFunc("POST /auth/login", r.Controller.Login)
	r.Mux.HandleFunc("POST /auth/register", r.Controller.Register)
	r.Mux.HandleFunc("POST /auth/verify/email/{email}", r.Controller.SendVerifyCode)
	r.Mux.HandleFunc("PUT /auth/verify/email/{email}/{code}", r.Controller.VerifyEmail)
	r.Mux.HandleFunc("PUT /auth/reset/password/verify/{email}/{code}/{new_password}", r.Controller.ResetPasswordVerify)
}
