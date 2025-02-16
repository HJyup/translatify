package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/HJyup/translatify-gateway/internal/gateway/chat"
	"io"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/internal/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	gateway chat.Gateway
}

func NewChatHandler(gateway chat.Gateway) *ChatHandler {
	return &ChatHandler{gateway: gateway}
}

func (h *ChatHandler) RegisterRoutes(router *mux.Router) {
	chatRouter := router.PathPrefix("/api/v1/chats").Subrouter()

	chatRouter.Handle("/{chatId}", utils.TokenAuthMiddleware(http.HandlerFunc(h.HandleChat))).Methods("GET")
	chatRouter.Handle("/{chatId}/messages", utils.TokenAuthMiddleware(http.HandlerFunc(h.HandleListMessages))).Methods("GET")
	chatRouter.Handle("", utils.TokenAuthMiddleware(http.HandlerFunc(h.HandleCreateChat))).Methods("POST")
	chatRouter.Handle("/{chatId}/messages", utils.TokenAuthMiddleware(http.HandlerFunc(h.HandleSendMessage))).Methods("POST")
	chatRouter.Handle("/{chatId}/messages/stream", utils.TokenAuthMiddleware(http.HandlerFunc(h.HandleStreamMessages))).Methods("GET")
}

func (h *ChatHandler) HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	var reqBody models.CreateChatRequest

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

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	req := &pb.CreateChatRequest{
		UsernameA:      reqBody.UserAId,
		UsernameB:      reqBody.UserBId,
		SourceLanguage: reqBody.SourceLanguage,
		TargetLanguage: reqBody.TargetLanguage,
	}

	resp, err := h.gateway.CreateChat(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ChatId := vars["ChatId"]
	if ChatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "ChatId is required")
		return
	}

	req := &pb.GetChatRequest{
		ChatId: ChatId,
	}

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	resp, err := h.gateway.GetChat(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId := vars["chatId"]
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "chatId is required")
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
		ChatId:           chatId,
		SenderUsername:   reqBody.FromUserID,
		ReceiverUsername: reqBody.ToUserID,
		Content:          reqBody.Content,
	}

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	resp, err := h.gateway.SendMessage(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) HandleListMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId := vars["chatId"]
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "chatId is required")
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
		ChatId:         chatId,
		SinceTimestamp: sinceTimestamp,
		Limit:          limit,
		PageToken:      pageToken,
	}

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	resp, err := h.gateway.ListMessages(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

// HandleGetMessage This handler is not used yet.
func (h *ChatHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	if messageId == "" {
		utils.WriteError(w, http.StatusBadRequest, "messageId is required")
		return
	}

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	req := &pb.GetMessageRequest{MessageId: messageId}
	resp, err := h.gateway.GetMessage(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) HandleStreamMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId := vars["chatId"]
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "chatId is required")
		return
	}

	req := &pb.StreamMessagesRequest{
		ChatId: chatId,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to upgrade connection: "+err.Error())
		return
	}
	defer conn.Close()

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	grpcStream, err := h.gateway.StreamMessages(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
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
