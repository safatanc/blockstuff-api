package payout

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/safatanc/blockstuff-api/internal/domain/item"
	"github.com/safatanc/blockstuff-api/internal/domain/transaction"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/safatanc/blockstuff-api/pkg/response"
	"github.com/safatanc/blockstuff-api/pkg/util"
)

type Controller struct {
	Service            *Service
	UserService        *user.Service
	ItemService        *item.Service
	TransactionService *transaction.Service
}

func NewController(service *Service, userService *user.Service, itemService *item.Service, transactionService *transaction.Service) *Controller {
	return &Controller{
		Service:            service,
		UserService:        userService,
		ItemService:        itemService,
		TransactionService: transactionService,
	}
}

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	status := r.URL.Query().Get("status")

	if !(claimsUser.Role == "ADMIN") {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	payouts := c.Service.FindAll(status)
	response.Success(w, payouts)
}

func (c *Controller) FindByID(w http.ResponseWriter, r *http.Request) {
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

	payout, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	response.Success(w, payout)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	claimsUsername := claims["username"].(string)
	claimsUser, err := c.UserService.FindByUsername(claimsUsername)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	var payout *Payout
	json.NewDecoder(r.Body).Decode(&payout)

	for _, payoutTransaction := range payout.PayoutTransactions {
		transaction, err := c.TransactionService.FindByID(payoutTransaction.TransactionID)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		for _, transactionItem := range transaction.TransactionItems {
			item, err := c.ItemService.FindByID(transactionItem.ItemID)
			if err != nil {
				response.Error(w, util.GetErrorStatusCode(err), err.Error())
				return
			}
			if !(claimsUser.Role == "ADMIN" || claimsUser.ID.String() == item.MinecraftServer.AuthorID) {
				response.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}
		}
	}

	payout, err = c.Service.Create(payout)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, payout)
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
	payout, err := c.Service.FindByID(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	json.NewDecoder(r.Body).Decode(&payout)

	payout, err = c.Service.Update(id, payout)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, payout)
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

	payout, err := c.Service.Delete(id)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, payout)
}

func (c *Controller) FindPayoutChannels(w http.ResponseWriter, r *http.Request) {
	payoutChannels, err := c.Service.FindPayoutChannels()
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}
	response.Success(w, payoutChannels)
}

func (c *Controller) GetPayoutChannel(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	claims := r.Context().Value("claims").(jwt.MapClaims)

	if claims["username"] != username {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := c.UserService.FindByUsername(username)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	userPayoutChannel, err := c.Service.GetPayoutChannel(user.ID.String())
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	response.Success(w, userPayoutChannel)
}

func (c *Controller) SetPayoutChannel(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	claims := r.Context().Value("claims").(jwt.MapClaims)

	if claims["username"] != username {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var userPayoutChannel *user.UserPayoutChannel
	err := json.NewDecoder(r.Body).Decode(&userPayoutChannel)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.UserService.FindByUsername(username)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	userPayoutChannel, err = c.Service.SetPayoutChannel(user.ID.String(), userPayoutChannel)
	if err != nil {
		response.Error(w, util.GetErrorStatusCode(err), err.Error())
		return
	}

	response.Success(w, userPayoutChannel)
}
