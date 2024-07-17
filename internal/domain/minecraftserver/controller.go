package minecraftserver

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

func (c *Controller) FindByIP(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	minecraftserver, err := c.Service.FindByIP(ip)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var minecraftserver *MinecraftServer
	json.NewDecoder(r.Body).Decode(&minecraftserver)

	minecraftserver, err := c.Service.Create(minecraftserver)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	minecraftserver, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	json.NewDecoder(r.Body).Decode(&minecraftserver)

	minecraftserver, err = c.Service.Update(id, minecraftserver)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	minecraftserver, err := c.Service.Delete(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}
