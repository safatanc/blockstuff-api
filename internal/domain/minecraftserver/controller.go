package minecraftserver

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

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
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 30
	}
	authorID := r.URL.Query().Get("author_id")

	minecraftservers := c.Service.FindAll(page, limit, authorID)
	response.Success(w, minecraftservers)
}

func (c *Controller) FindBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	minecraftserver, err := c.Service.FindBySlug(slug, false)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}

func (c *Controller) FindBySlugDetail(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	slug := r.PathValue("slug")
	minecraftserver, err := c.Service.FindBySlug(slug, true)
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

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
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

func (c *Controller) UpdateLogo(w http.ResponseWriter, r *http.Request) {
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

	minecraftserver, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == minecraftserver.AuthorID) {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	minecraftserver, err = c.Service.UpdateLogo(id, uploadedImage, imageHeader)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, minecraftserver)
}
