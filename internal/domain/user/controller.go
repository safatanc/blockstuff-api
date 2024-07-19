package user

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/safatanc/blockstuff-api/pkg/response"
	"github.com/safatanc/blockstuff-api/pkg/util"
)

type Controller struct {
	Service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{
		Service: service,
	}
}

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	users := c.Service.FindAll()
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.Service.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if claimsUser.Role != "ADMIN" {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	response.Success(w, users)
}

func (c *Controller) FindByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	claims := r.Context().Value("claims").(jwt.MapClaims)

	if claims["username"] != username {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := c.Service.FindByUsername(username)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	response.Success(w, user)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.Service.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if claimsUser.Role != "ADMIN" {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var user *User
	json.NewDecoder(r.Body).Decode(&user)

	user, err = c.Service.Create(user)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, user)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.Service.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	id := r.PathValue("id")
	user, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUsername == user.Username) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	json.NewDecoder(r.Body).Decode(&user)

	if user.Role == "ADMIN" {
		if claimsUser.Role != "ADMIN" {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
	}

	user, err = c.Service.Update(id, user)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, user)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.Service.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	id := r.PathValue("id")
	user, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUsername == user.Username) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err = c.Service.Delete(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, user)
}
