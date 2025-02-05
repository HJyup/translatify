package handlers

import (
	"encoding/json"
	"github.com/HJyup/translatify-gateway/internal/gateway"
	"github.com/HJyup/translatify-gateway/internal/models"
	"io"
	"net/http"
	"time"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	router.HandleFunc("/api/v1/messages/list", h.HandleListMessages).Methods("POST")
	router.HandleFunc("/api/v1/messages", h.HandleSendMessage).Methods("POST")
	router.HandleFunc("/api/v1/messages/stream", h.HandleStreamMessages).Methods("GET")

	router.HandleFunc("/api/v1/messages/{messageID}", h.HandleGetMessage).Methods("GET")
}

func (h *ChatHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
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
		FromUserId:     reqBody.FromUserID,
		ToUserId:       reqBody.ToUserID,
		Content:        reqBody.Content,
		SourceLanguage: reqBody.SourceLang,
		TargetLanguage: reqBody.TargetLang,
	}

	res, err := h.gateway.AddMessage(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, res)
}

func (h *ChatHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["messageID"]

	res, err := h.gateway.GetMessage(r.Context(), &pb.GetMessageRequest{
		MessageId: messageID,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (h *ChatHandler) HandleListMessages(w http.ResponseWriter, r *http.Request) {
	var reqBody models.ListMessagesRequest
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

	req := &pb.ListMessagesRequest{
		UserId:              reqBody.UserID,
		CorrespondentUserId: reqBody.CorrespondentUserID,
		SinceTimestamp:      reqBody.SinceTimestamp,
	}

	res, err := h.gateway.ListMessages(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (h *ChatHandler) HandleStreamMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("userId")
	correspondentID := query.Get("correspondentId")

	if userID == "" || correspondentID == "" {
		utils.WriteError(w, http.StatusBadRequest, "missing query parameters: userId and correspondentId are required")
		return
	}

	now := time.Now().Unix()

	streamReq := &pb.StreamMessagesRequest{
		UserId:              userID,
		CorrespondentUserId: correspondentID,
		SinceTimestamp:      now,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to upgrade connection: "+err.Error())
		return
	}
	defer conn.Close()

	grpcStream, err := h.gateway.StreamMessages(r.Context(), streamReq)
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
