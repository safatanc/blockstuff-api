package user

import (
	"encoding/json"
	"net/http"

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
	response.Success(w, users)
}

func (c *Controller) FindByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	user, err := c.Service.FindByUsername(username)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, user)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var user *User
	json.NewDecoder(r.Body).Decode(&user)

	user, err := c.Service.Create(user)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, user)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	json.NewDecoder(r.Body).Decode(&user)

	user, err = c.Service.Update(id, user)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, user)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := c.Service.Delete(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, user)
}
