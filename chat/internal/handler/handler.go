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

func (h *GrpcHandler) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	userNameA := req.GetUsernameA()
	userNameB := req.GetUsernameB()
	sourceLang := req.GetSourceLanguage()
	targetLang := req.GetTargetLanguage()

	if userNameA == "" || userNameB == "" || sourceLang == "" || targetLang == "" {
		return nil, status.Error(codes.InvalidArgument, "username_a username_b, source_lang, and target_lang must be provided")
	}

	convID, err := h.service.CreateChat(userNameA, userNameB, sourceLang, targetLang)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Chat: %v", err)
	}

	return &pb.CreateChatResponse{
		Success: true,
		ChatId:  convID,
		Error:   "",
	}, nil
}

func (h *GrpcHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	conID := req.GetChatId()
	senderUserName := req.GetSenderUsername()
	receiverUserName := req.GetReceiverUsername()
	content := req.GetContent()

	if conID == "" || senderUserName == "" || content == "" || receiverUserName == "" {
		return nil, status.Error(codes.InvalidArgument, "Chat_id, sender_username, receiver_username, and content must be provided")
	}

	conv, err := h.service.GetChat(conID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get Chat: %v", err)
	}

	messageID, err := h.service.SendMessage(conID, senderUserName, receiverUserName, content)
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
	convID := req.GetChatId()

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
	if req.GetChatId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Chat_id must be provided")
	}

	var since *timestamppb.Timestamp
	if req.GetSinceTimestamp() > 0 {
		since = timestamppb.New(time.Unix(req.GetSinceTimestamp(), 0))
	}

	messages, pageToken, err := h.service.ListMessages(req.GetChatId(), since, int(req.GetLimit()), req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	resp := &pb.ListMessagesResponse{
		Messages:      messages,
		NextPageToken: pageToken,
	}
	return resp, nil
}

func (h *GrpcHandler) GetChat(ctx context.Context, req *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	conv, err := h.service.GetChat(req.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get Chat: %v", err)
	}

	return &pb.GetChatResponse{Chat: conv}, nil
}

func (h *GrpcHandler) ListChats(ctx context.Context, req *pb.ListChatsRequest) (*pb.ListChatsResponse, error) {
	userName := req.GetUsername()
	if userName == "" {
		return nil, status.Error(codes.InvalidArgument, "username must be provided")
	}

	conv, err := h.service.ListChats(userName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list Chats: %v", err)
	}

	return &pb.ListChatsResponse{Chats: conv}, nil
}
