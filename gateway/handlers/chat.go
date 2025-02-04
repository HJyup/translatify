package handlers

import (
	common "github.com/HJyup/translatify-common"
	pb "github.com/HJyup/translatify-common/api"
	"net/http"
)

type ChatHandler struct {
	chatClient pb.ChatServiceClient
}

func NewChatHandler(chatClient pb.ChatServiceClient) *ChatHandler {
	return &ChatHandler{chatClient}
}

func (h *ChatHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/messages/{messageID}", h.HandleGetMessage)
}

func (h *ChatHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	messageID := r.PathValue("messageID")

	res, err := h.chatClient.GetMessage(r.Context(), &pb.GetMessageRequest{
		MessageId: messageID,
	})
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	common.WriteJSON(w, 200, res)
}
