package callback

import (
	"encoding/json"
	"net/http"
	"os"

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

func (c *Controller) XenditCallback(w http.ResponseWriter, r *http.Request) {
	requestCallbackToken := r.Header.Get("X-CALLBACK-TOKEN")
	if requestCallbackToken == "" {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if requestCallbackToken != os.Getenv("XENDIT_CALLBACK_KEY") {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var xenditPayload *XenditPayload
	json.NewDecoder(r.Body).Decode(&xenditPayload)

	err := c.Service.XenditCallback(xenditPayload)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, nil)
}
