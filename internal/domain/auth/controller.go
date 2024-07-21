package auth

import (
	"encoding/json"
	"net/http"

	"github.com/safatanc/blockstuff-api/internal/domain/user"
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

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var user *user.User
	json.NewDecoder(r.Body).Decode(&user)

	auth, err := c.Service.VerifyUser(user.Username, user.Password)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, auth)
}