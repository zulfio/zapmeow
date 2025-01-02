package handler

import (
	"net/http"
	"zapmeow/api/response"
	"zapmeow/api/service"
	"zapmeow/pkg/whatsapp"

	"github.com/gin-gonic/gin"
)

type getContactsResponse struct {
	Contacts []whatsapp.ContactInfo `json:"contacts"`
}

type getContactsHandler struct {
	whatsAppService service.WhatsAppService
}

func NewGetContactsHandler(
	whatsAppService service.WhatsAppService,
) *getContactsHandler {
	return &getContactsHandler{
		whatsAppService: whatsAppService,
	}
}

// Get WhatsApp Contacts
//
//	@Summary		Get WhatsApp Contacts
//	@Description	Retrieves all contacts for the specified WhatsApp instance.
//	@Tags			WhatsApp Contacts
//	@Param			instanceId	path	string	true	"Instance ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	getContactsResponse	"List of contacts"
//	@Failure		401	{object}	response.Error		"Unauthorized"
//	@Failure		500	{object}	response.Error		"Internal Server Error"
//	@Router			/{instanceId}/contacts [get]
func (h *getContactsHandler) Handler(c *gin.Context) {
	instanceID := c.Param("instanceId")
	instance, err := h.whatsAppService.GetInstance(instanceID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !h.whatsAppService.IsAuthenticated(instance) {
		response.ErrorResponse(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	contacts, err := h.whatsAppService.GetContacts(instance)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Response(c, http.StatusOK, getContactsResponse{
		Contacts: contacts,
	})
}