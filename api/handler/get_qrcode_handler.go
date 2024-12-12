package handler

import (
	"net/http"
	"zapmeow/api/helper"
	"zapmeow/api/response"
	"zapmeow/api/service"
	"zapmeow/pkg/whatsapp"
	"zapmeow/pkg/zapmeow"

	"github.com/gin-gonic/gin"
)

type getQrCodeResponse struct {
	QrCode string             `json:"qrcode"`
	Status string             `json:"status"`
	Info   *whatsapp.ContactInfo `json:"info,omitempty"` // Add this field
}

type getQrCodeHandler struct {
	app             *zapmeow.ZapMeow
	whatsAppService service.WhatsAppService
	messageService  service.MessageService
	accountService  service.AccountService
}

func NewGetQrCodeHandler(
	app *zapmeow.ZapMeow,
	whatsAppService service.WhatsAppService,
	messageService service.MessageService,
	accountService service.AccountService,
) *getQrCodeHandler {
	return &getQrCodeHandler{
		app:             app,
		whatsAppService: whatsAppService,
		messageService:  messageService,
		accountService:  accountService,
	}
}

// Get QR Code for WhatsApp Login
//
//	@Summary		Get WhatsApp QR Code
//	@Description	Returns a QR code to initiate WhatsApp login.
//	@Tags			WhatsApp Login
//	@Param			instanceId	path	string	true	"Instance ID"
//	@Produce		json
//	@Success		200	{object}	getQrCodeResponse	"QR Code and Profile Info"
//	@Router			/{instanceId}/qrcode [get]
func (h *getQrCodeHandler) Handler(c *gin.Context) {
	instanceID := c.Param("instanceId")
	instance, err := h.whatsAppService.GetInstance(instanceID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.app.Mutex.Lock()
	defer h.app.Mutex.Unlock()
	account, err := h.accountService.GetAccountByInstanceID(instanceID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if account == nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Account not found")
		return
	}

	var responseStatus string
	var profileInfo *whatsapp.ContactInfo

	if account.Status == "CONNECTED" {
		responseStatus = "CONNECTED"
		
		// Get profile info when connected
		jid, ok := helper.MakeJID(account.User)
		if !ok {
			response.ErrorResponse(c, http.StatusInternalServerError, "Invalid user JID")
			return
		}

		info, err := h.whatsAppService.GetContactInfo(instance, jid)
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		profileInfo = info
	} else {
		responseStatus = "SCAN_QR_CODE"
	}

	response.Response(c, http.StatusOK, getQrCodeResponse{
		QrCode: account.QrCode,
		Status: responseStatus,
		Info:   profileInfo,
	})
}