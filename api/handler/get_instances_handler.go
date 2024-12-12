package handler

import (
	"net/http"
	"zapmeow/api/response"
	"zapmeow/api/service"

	"github.com/gin-gonic/gin"
)

type getInstancesResponse struct {
	Instances []instanceResponse `json:"instances"`
}

type instanceResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type getInstancesHandler struct {
	accountService service.AccountService
}

func NewGetInstancesHandler(
	accountService service.AccountService,
) *getInstancesHandler {
	return &getInstancesHandler{
		accountService: accountService,
	}
}

// Get WhatsApp Instances
//
//	@Summary		Get WhatsApp Instances
//	@Description	Returns all WhatsApp instances.
//	@Tags			WhatsApp Instance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	getInstancesResponse	"List of instances"
//	@Router			/instances [get]
func (h *getInstancesHandler) Handler(c *gin.Context) {
	accounts, err := h.accountService.GetAllAccounts()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var instances []instanceResponse
	for _, account := range accounts {
		instances = append(instances, instanceResponse{
			ID:     account.InstanceID,
			Status: account.Status,
		})
	}

	response.Response(c, http.StatusOK, getInstancesResponse{
		Instances: instances,
	})
}