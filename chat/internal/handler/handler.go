package handler

import (
	"context"
	"encoding/json"
	"github.com/HJyup/translatify-common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"

	"github.com/HJyup/translatify-chat/internal/models"
	pb "github.com/HJyup/translatify-common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcHandler struct {
	pb.UnimplementedChatServiceServer

	service models.ChatService
	channel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service models.ChatService, channel *amqp.Channel) {
	handler := &GrpcHandler{
		service: service,
		channel: channel,
	}

	pb.RegisterChatServiceServer(grpcServer, handler)
}

func (h *GrpcHandler) CreateConversation(ctx context.Context, req *pb.CreateConversationRequest) (*pb.CreateConversationResponse, error) {
	userAID := req.GetUserAId()
	userBID := req.GetUserBId()
	sourceLang := req.GetSourceLanguage()
	targetLang := req.GetTargetLanguage()

	if userAID == "" || userBID == "" || sourceLang == "" || targetLang == "" {
		return nil, status.Error(codes.InvalidArgument, "user_a_id, user_b_id, source_lang, and target_lang must be provided")
	}

	convID, err := h.service.CreateConversation(userAID, userBID, sourceLang, targetLang)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create conversation: %v", err)
	}

	return &pb.CreateConversationResponse{
		Success:        true,
		ConversationId: convID,
		Error:          "",
	}, nil
}

func (h *GrpcHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	conID := req.GetConversationId()
	senderID := req.GetSenderId()
	receiverID := req.GetReceiverId()
	content := req.GetContent()

	if conID == "" || senderID == "" || content == "" || receiverID == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id, sender_id, receiver_id, and content must be provided")
	}

	conv, err := h.service.GetConversation(conID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get conversation: %v", err)
	}

	messageID, err := h.service.SendMessage(conID, senderID, receiverID, content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	q, err := h.channel.QueueDeclare(broker.MessageSentEvent, true, false, false, false, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to declare queue: %v", err)
	}

	msgData := map[string]interface{}{
		"sourceLang": conv.GetTargetLanguage(),
		"targetLang": conv.GetSourceLanguage(),
		"messageID":  messageID,
		"content":    content,
	}

	body, err := json.Marshal(msgData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal message to JSON: %v", err)
	}

	err = h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	}

	return &pb.SendMessageResponse{MessageId: messageID}, nil
}

func (h *GrpcHandler) StreamMessages(req *pb.StreamMessagesRequest, stream pb.ChatService_StreamMessagesServer) error {
	convID := req.GetConversationId()

	msgCh, err := h.service.StreamMessages(stream.Context(), convID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to start message stream: %v", err)
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg, ok := <-msgCh:
			if !ok {
				return nil
			}
			if err = stream.Send(msg); err != nil {
				return status.Errorf(codes.Internal, "failed to send message: %v", err)
			}
		}
	}
}

func (h *GrpcHandler) GetMessage(ctx context.Context, req *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	msg, err := h.service.GetMessage(req.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get message: %v", err)
	}

	return &pb.GetMessageResponse{Message: msg}, nil
}

func (h *GrpcHandler) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	if req.GetConversationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id must be provided")
	}

	var since *timestamppb.Timestamp
	if req.GetSinceTimestamp() > 0 {
		since = timestamppb.New(time.Unix(req.GetSinceTimestamp(), 0))
	}

	messages, pageToken, err := h.service.ListMessages(req.GetConversationId(), since, int(req.GetLimit()), req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	resp := &pb.ListMessagesResponse{
		Messages:      messages,
		NextPageToken: pageToken,
	}
	return resp, nil
}

func (h *GrpcHandler) GetConversation(ctx context.Context, req *pb.GetConversationRequest) (*pb.GetConversationResponse, error) {
	conv, err := h.service.GetConversation(req.GetConversationId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get conversation: %v", err)
	}

	return &pb.GetConversationResponse{Conversation: conv}, nil
}

func (h *GrpcHandler) ListConversations(ctx context.Context, req *pb.ListConversationsRequest) (*pb.ListConversationsResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id must be provided")
	}

	conv, err := h.service.ListConversations(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list conversations: %v", err)
	}

	return &pb.ListConversationsResponse{Conversations: conv}, nil
}
