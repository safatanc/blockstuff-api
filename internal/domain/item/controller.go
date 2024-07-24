package item

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

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
	minecraftServerID := r.PathValue("minecraft_server_id")
	items := c.Service.FindAll(minecraftServerID)
	response.Success(w, items)
}

func (c *Controller) FindBySlug(w http.ResponseWriter, r *http.Request) {
	minecraftServerID := r.PathValue("minecraft_server_id")
	slug := r.PathValue("slug")

	item, err := c.Service.FindBySlug(minecraftServerID, slug)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	response.Success(w, item)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	minecraftServerID := r.PathValue("minecraft_server_id")

	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	minecraftserver, err := c.MinecraftServerService.FindByID(minecraftServerID)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var item *Item
	json.NewDecoder(r.Body).Decode(&item)
	item.MinecraftServerID = &minecraftServerID

	item, err = c.Service.Create(item)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, item)
}

func (c *Controller) AddImage(w http.ResponseWriter, r *http.Request) {
	minecraftServerID := r.PathValue("minecraft_server_id")
	id := r.PathValue("id")

	maxUploadSizeMB, err := strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE_MB"))
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	err = r.ParseMultipartForm(int64(maxUploadSizeMB))
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	uploadedImage, imageHeader, err := r.FormFile("image")
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	defer uploadedImage.Close()

	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	minecraftserver, err := c.MinecraftServerService.FindByID(minecraftServerID)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemImage, err := c.Service.AddImage(id, uploadedImage, imageHeader)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, itemImage)
}

func (c *Controller) AddAction(w http.ResponseWriter, r *http.Request) {
	minecraftServerID := r.PathValue("minecraft_server_id")
	id := r.PathValue("id")

	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	minecraftserver, err := c.MinecraftServerService.FindByID(minecraftServerID)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var itemAction *ItemAction
	json.NewDecoder(r.Body).Decode(&itemAction)
	itemAction.ItemID = id

	itemAction, err = c.Service.AddAction(itemAction)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, itemAction)
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
	item, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	minecraftserver, err := c.MinecraftServerService.FindByID(*item.MinecraftServerID)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	json.NewDecoder(r.Body).Decode(&item)

	item, err = c.Service.Update(id, item)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, item)
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
	item, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	minecraftserver, err := c.MinecraftServerService.FindByID(*item.MinecraftServerID)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	item, err = c.Service.Delete(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, item)
}
