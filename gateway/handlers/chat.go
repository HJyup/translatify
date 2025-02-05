package handlers

import (
	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/gateway"
	"net/http"
)

type ChatHandler struct {
	gateway gateway.ChatGateway
}

func NewChatHandler(gateway gateway.ChatGateway) *ChatHandler {
	return &ChatHandler{gateway}
}

func (h *ChatHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/messages/{messageID}", h.HandleGetMessage)
}

func (h *ChatHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	messageID := r.PathValue("messageID")

	res, err := h.gateway.GetMessage(r.Context(), &pb.GetMessageRequest{
		MessageId: messageID,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, 200, res)
}
