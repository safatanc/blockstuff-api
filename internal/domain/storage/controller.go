package storage

import (
	"bytes"
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

func (c *Controller) Find(w http.ResponseWriter, r *http.Request) {
	objectName := r.PathValue("object_name")
	object, err := c.Service.Find(objectName)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)
	w.Write(buf.Bytes())
}
