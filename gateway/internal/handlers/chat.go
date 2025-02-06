package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/internal/gateway"
	"github.com/HJyup/translatify-gateway/internal/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	gateway gateway.ChatGateway
}

func NewChatHandler(gateway gateway.ChatGateway) *ChatHandler {
	return &ChatHandler{gateway: gateway}
}

func (h *ChatHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/conversations", h.HandleCreateConversation).Methods("POST")
	router.HandleFunc("/api/v1/conversations/{conversationId}/messages", h.HandleSendMessage).Methods("POST")
	router.HandleFunc("/api/v1/conversations/{conversationId}/messages", h.HandleListMessages).Methods("GET")
	router.HandleFunc("/api/v1/messages/{messageId}", h.HandleGetMessage).Methods("GET")
	router.HandleFunc("/api/v1/conversations/{conversationId}/messages/stream", h.HandleStreamMessages).Methods("GET")
}

func (h *ChatHandler) HandleCreateConversation(w http.ResponseWriter, r *http.Request) {
	var reqBody models.CreateConversationRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &reqBody); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	req := &pb.CreateConversationRequest{
		UserAId:        reqBody.UserAId,
		UserBId:        reqBody.UserBId,
		SourceLanguage: reqBody.SourceLanguage,
		TargetLanguage: reqBody.TargetLanguage,
	}

	resp, err := h.gateway.CreateConversation(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationId := vars["conversationId"]
	if conversationId == "" {
		utils.WriteError(w, http.StatusBadRequest, "conversationId is required")
		return
	}

	var reqBody models.SendMessageRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &reqBody); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	req := &pb.SendMessageRequest{
		ConversationId: conversationId,
		SenderId:       reqBody.FromUserID,
		ReceiverId:     reqBody.ToUserID,
		Content:        reqBody.Content,
	}

	resp, err := h.gateway.SendMessage(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) HandleListMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationId := vars["conversationId"]
	if conversationId == "" {
		utils.WriteError(w, http.StatusBadRequest, "conversationId is required")
		return
	}

	q := r.URL.Query()
	sinceStr := q.Get("sinceTimestamp")
	limitStr := q.Get("limit")
	pageToken := q.Get("pageToken")

	var sinceTimestamp int64
	if sinceStr != "" {
		var err error
		sinceTimestamp, err = strconv.ParseInt(sinceStr, 10, 64)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid sinceTimestamp")
			return
		}
	}

	var limit int32 = 50
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = int32(l)
	}

	req := &pb.ListMessagesRequest{
		ConversationId: conversationId,
		SinceTimestamp: sinceTimestamp,
		Limit:          limit,
		PageToken:      pageToken,
	}

	resp, err := h.gateway.ListMessages(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	if messageId == "" {
		utils.WriteError(w, http.StatusBadRequest, "messageId is required")
		return
	}

	req := &pb.GetMessageRequest{MessageId: messageId}
	resp, err := h.gateway.GetMessage(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) HandleStreamMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationId := vars["conversationId"]
	if conversationId == "" {
		utils.WriteError(w, http.StatusBadRequest, "conversationId is required")
		return
	}

	req := &pb.StreamMessagesRequest{
		ConversationId: conversationId,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to upgrade connection: "+err.Error())
		return
	}
	defer conn.Close()

	grpcStream, err := h.gateway.StreamMessages(r.Context(), req)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": err.Error()})
		return
	}

	for {
		msg, err := grpcStream.Recv()
		if err != nil {
			conn.WriteJSON(map[string]string{"error": err.Error()})
			return
		}
		if err := conn.WriteJSON(msg); err != nil {
			return
		}
	}
}
