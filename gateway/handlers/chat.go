package handlers

import (
	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/gateway"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	gateway gateway.ChatGateway
}

func NewChatHandler(gateway gateway.ChatGateway) *ChatHandler {
	return &ChatHandler{gateway}
}

func (h *ChatHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/messages/{messageID}", h.HandleGetMessage)
	mux.HandleFunc("GET /api/v1/messages/{userID}/{correspondentID}/{sinceTimestamp}", h.HandleGetListOfMessages)
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

func (h *ChatHandler) HandleGetListOfMessages(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	correspondentID := r.PathValue("correspondentID")
	sinceTimestampStr := r.PathValue("sinceTimestamp")

	var sinceTimestamp int64
	if sinceTimestampStr != "" {
		var err error
		sinceTimestamp, err = strconv.ParseInt(sinceTimestampStr, 10, 64)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid sinceTimestamp")
			return
		}
	}

	req := &pb.ListMessagesRequest{
		UserId:              userID,
		CorrespondentUserId: correspondentID,
		SinceTimestamp:      sinceTimestamp,
	}

	res, err := h.gateway.ListMessages(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, res)
}
