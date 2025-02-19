package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/internal/gateway/chat"
	"github.com/HJyup/translatify-gateway/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

func extractUsername(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid authorization header")
	}
	claims, err := utils.ParseToken(strings.TrimSpace(parts[1]))
	if err != nil {
		return "", err
	}
	return claims.UserName, nil
}

// HandleCreateChat godoc
// @Summary Create Chat
// @Description Create a new chat between two users.
// @Tags chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param chat body models.CreateChatRequest true "Chat information"
// @Success 200 {object} api.CreateChatResponse "Chat created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/chats [post]
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
	tokenUsername, err := extractUsername(r)
	if err != nil || (tokenUsername != reqBody.UserNameA && tokenUsername != reqBody.UserNameB) {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	req := &api.CreateChatRequest{
		UsernameA:      reqBody.UserNameA,
		UsernameB:      reqBody.UserNameB,
		SourceLanguage: reqBody.SourceLanguage,
		TargetLanguage: reqBody.TargetLanguage,
	}
	resp, err := h.gateway.CreateChat(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

// HandleChat godoc
// @Summary Get Chat
// @Description Get details of a specific chat.
// @Tags chats
// @Security BearerAuth
// @Produce json
// @Param chatId path string true "Chat ID"
// @Success 200 {object} api.GetChatResponse "Chat details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/chats/{chatId} [get]
func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId := vars["chatId"]
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "chatId is required")
		return
	}
	getReq := &api.GetChatRequest{ChatId: chatId}
	chatResp, err := h.gateway.GetChat(r.Context(), getReq)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tokenUsername, err := extractUsername(r)
	if err != nil || (tokenUsername != chatResp.Chat.UsernameA && tokenUsername != chatResp.Chat.UsernameB) {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	utils.WriteJSON(w, http.StatusOK, chatResp)
}

// HandleListMessages godoc
// @Summary List Messages
// @Description Get a list of messages for a specific chat.
// @Tags chats
// @Security BearerAuth
// @Produce json
// @Param chatId path string true "Chat ID"
// @Param sinceTimestamp query int false "Since timestamp (Unix epoch in seconds)"
// @Param limit query int false "Maximum number of messages to return"
// @Param pageToken query string false "Token for pagination"
// @Success 200 {object} api.ListMessagesResponse "List of messages"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/chats/{chatId}/messages [get]
func (h *ChatHandler) HandleListMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId := vars["chatId"]
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "chatId is required")
		return
	}
	getReq := &api.GetChatRequest{ChatId: chatId}
	chatResp, err := h.gateway.GetChat(r.Context(), getReq)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tokenUsername, err := extractUsername(r)
	if err != nil || (tokenUsername != chatResp.Chat.UsernameA && tokenUsername != chatResp.Chat.UsernameB) {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
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
	req := &api.ListMessagesRequest{
		ChatId:         chatId,
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

// HandleStreamMessages godoc
// @Summary Stream Messages
// @Description Open a websocket connection to stream messages for a specific chat.
// @Tags chats
// @Security BearerAuth
// @Produce json
// @Param chatId path string true "Chat ID"
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/chats/{chatId}/messages/stream [get]
func (h *ChatHandler) HandleStreamMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId := vars["chatId"]
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "chatId is required")
		return
	}
	getReq := &api.GetChatRequest{ChatId: chatId}
	chatResp, err := h.gateway.GetChat(r.Context(), getReq)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tokenUsername, err := extractUsername(r)
	if err != nil || (tokenUsername != chatResp.Chat.UsernameA && tokenUsername != chatResp.Chat.UsernameB) {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to upgrade connection: "+err.Error())
		return
	}
	defer conn.Close()
	req := &api.StreamMessagesRequest{
		ChatId: chatId,
	}
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

// HandleSendMessage godoc
// @Summary Send Message
// @Description Send a message in a chat.
// @Tags chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param chatId path string true "Chat ID"
// @Param message body models.SendMessageRequest true "Message information"
// @Success 200 {object} api.SendMessageResponse "Message sent successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/chats/{chatId}/messages [post]
func (h *ChatHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	vars := mux.Vars(r)
	chatId := vars["chatId"]
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "chatId is required")
		return
	}
	getReq := &api.GetChatRequest{ChatId: chatId}
	chatResp, err := h.gateway.GetChat(ctx, getReq)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tokenUsername, err := extractUsername(r)
	if err != nil || (tokenUsername != chatResp.Chat.UsernameA && tokenUsername != chatResp.Chat.UsernameB) {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
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
	if tokenUsername != reqBody.FromUserName {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized sender")
		return
	}
	req := &api.SendMessageRequest{
		ChatId:           chatId,
		SenderUsername:   reqBody.FromUserName,
		ReceiverUsername: reqBody.ToUserName,
		Content:          reqBody.Content,
	}
	resp, err := h.gateway.SendMessage(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}
