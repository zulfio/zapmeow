package handler

import (
	"net/http"
	"zapmeow/api/helper"
	"zapmeow/api/model"
	"zapmeow/api/response"
	"zapmeow/api/service"

	"github.com/gin-gonic/gin"
	"github.com/vincent-petithory/dataurl"
)

type sendVideoMessageBody struct {
	Phone  string `json:"phone"`
	Base64 string `json:"base64"`
}

type sendVideoMessageResponse struct {
	Message response.Message `json:"message"`
}

type sendVideoMessageHandler struct {
	whatsAppService service.WhatsAppService
	messageService  service.MessageService
}

func NewSendVideoMessageHandler(
	whatsAppService service.WhatsAppService,
	messageService service.MessageService,
) *sendVideoMessageHandler {
	return &sendVideoMessageHandler{
		whatsAppService: whatsAppService,
		messageService:  messageService,
	}
}

// Send Video Message on WhatsApp
//
//	@Summary		Send Video Message on WhatsApp
//	@Description	Sends a video message on WhatsApp using the specified instance.
//	@Tags			WhatsApp Chat
//	@Param			instanceId	path	string					true	"Instance ID"
//	@Param			data		body	sendVideoMessageBody	true	"Video message body"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	sendVideoMessageResponse	"Message Send Response"
//	@Router			/{instanceId}/chat/send/video [post]
func (h *sendVideoMessageHandler) Handler(c *gin.Context) {
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

	var body sendVideoMessageBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Error trying to validate infos.")
		return
	}

	jid, ok := helper.MakeJID(body.Phone)
	if !ok {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid phone")
		return
	}

	mimetype, err := helper.GetMimeTypeFromDataURI(body.Base64)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	videoURL, err := dataurl.DecodeString(body.Base64)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.whatsAppService.SendVideoMessage(instance, jid, videoURL, mimetype)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	path, err := helper.SaveMedia(
		instanceID,
		resp.ID,
		videoURL.Data,
		mimetype,
	)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	message := model.Message{
		FromMe:     true,
		ChatJID:    jid.User,
		SenderJID:  resp.Sender.User,
		InstanceID: instanceID,
		Timestamp:  resp.Timestamp,
		MessageID:  resp.ID,
		MediaType:  "video",
		MediaPath:  path,
	}

	err = h.messageService.CreateMessage(&message)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Response(c, http.StatusOK, sendVideoMessageResponse{
		Message: response.NewMessageResponse(message),
	})
}