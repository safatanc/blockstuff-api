package minecraftserver

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/safatanc/blockstuff-api/pkg/response"
	"github.com/safatanc/blockstuff-api/pkg/util"
)

type Controller struct {
	Service     *Service
	UserService *user.Service
}

func NewController(service *Service, userService *user.Service) *Controller {
	return &Controller{
		Service:     service,
		UserService: userService,
	}
}

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	minecraftservers := c.Service.FindAll()
	response.Success(w, minecraftservers)
}

func (c *Controller) FindByIP(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	minecraftserver, err := c.Service.FindByIP(ip, false)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}

func (c *Controller) FindByIPDetail(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	ip := r.PathValue("ip")
	minecraftserver, err := c.Service.FindByIP(ip, true)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if claimsUser.ID.String() != minecraftserver.AuthorID {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
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
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	id := r.PathValue("id")
	minecraftserver, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if claimsUser.ID.String() != minecraftserver.AuthorID {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
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
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	id := r.PathValue("id")
	minecraftserver, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	minecraftserver, err = c.Service.Delete(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}

func (c *Controller) UpdateRcon(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	id := r.PathValue("id")
	minecraftserver, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var rcon *MinecraftServerRcon
	json.NewDecoder(r.Body).Decode(&rcon)

	rcon.MinecraftServerID = id

	minecraftserver, err = c.Service.UpdateRcon(rcon)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}
