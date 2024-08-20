package transaction

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/safatanc/blockstuff-api/internal/domain/minecraftserver"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/safatanc/blockstuff-api/pkg/response"
	"github.com/safatanc/blockstuff-api/pkg/util"
)

type Controller struct {
	Service                *Service
	UserService            *user.Service
	MinecraftServerService *minecraftserver.Service
}

func NewController(service *Service, userService *user.Service, minecraftServerService *minecraftserver.Service) *Controller {
	return &Controller{
		Service:                service,
		UserService:            userService,
		MinecraftServerService: minecraftServerService,
	}
}

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	serverId := r.URL.Query().Get("server_id")

	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if serverId != "" {
		minecraftserver, err := c.MinecraftServerService.FindByID(serverId)
		if err != nil {
			response.Error(w, util.GetErrorStatusCode(err), err.Error())
			return
		}

		if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
	} else {
		if !(claimsUser.Role == "ADMIN") {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
	}

	transactions := c.Service.FindAll(serverId)
	response.Success(w, transactions)
}

func (c *Controller) FindByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	item, err := c.Service.FindByCode(code)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	response.Success(w, item)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var transaction *Transaction
	json.NewDecoder(r.Body).Decode(&transaction)

	transaction, err := c.Service.Create(transaction)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, transaction)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN") {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.PathValue("id")
	transaction, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	json.NewDecoder(r.Body).Decode(&transaction)

	transaction, err = c.Service.Update(id, transaction)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, transaction)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN") {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.PathValue("id")
	_, err = c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	transaction, err := c.Service.Delete(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, transaction)
}
