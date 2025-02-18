package handler

import (
	"context"
	"github.com/HJyup/translatify-common/broker"
	"github.com/go-jose/go-jose/v3/json"
	"go.opentelemetry.io/otel"
	"time"

	"github.com/HJyup/translatify-chat/internal/models"
	pb "github.com/HJyup/translatify-common/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func chatFromModel(chat *models.Chat) *pb.Chat {
	if chat == nil {
		return nil
	}
	return &pb.Chat{
		ChatId:         chat.ChatID,
		UsernameA:      chat.UsernameA,
		UsernameB:      chat.UsernameB,
		CreatedAt:      chat.CreatedAt.Unix(),
		SourceLanguage: chat.SourceLang,
		TargetLanguage: chat.TargetLang,
	}
}

func chatMessageFromModel(msg *models.ChatMessage) *pb.ChatMessage {
	if msg == nil {
		return nil
	}
	return &pb.ChatMessage{
		MessageId:         msg.MessageID,
		ChatId:            msg.ChatID,
		SenderUsername:    msg.SenderUsername,
		ReceiverUsername:  msg.ReceiverUsername,
		Content:           msg.Content,
		TranslatedContent: msg.TranslatedContent,
		Timestamp:         msg.Timestamp.Unix(),
	}
}

func (h *GrpcHandler) CreateChat(_ context.Context, req *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	userNameA := req.GetUsernameA()
	userNameB := req.GetUsernameB()
	sourceLang := req.GetSourceLanguage()
	targetLang := req.GetTargetLanguage()

	if userNameA == "" || userNameB == "" || sourceLang == "" || targetLang == "" {
		return nil, status.Error(codes.InvalidArgument, "username_a, username_b, source_lang, and target_lang must be provided")
	}

	chatID, err := h.service.CreateChat(userNameA, userNameB, sourceLang, targetLang)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Chat: %v", err)
	}

	return &pb.CreateChatResponse{
		Success: true,
		ChatId:  chatID,
		Error:   "",
	}, nil
}

func (h *GrpcHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	ctx, span := otel.Tracer("chat-handler").Start(ctx, "SendMessage")
	defer span.End()

	chatID := req.GetChatId()
	senderUsername := req.GetSenderUsername()
	receiverUsername := req.GetReceiverUsername()
	content := req.GetContent()

	if chatID == "" || senderUsername == "" || receiverUsername == "" || content == "" {
		return nil, status.Error(codes.InvalidArgument, "chat_id, sender_username, receiver_username, and content must be provided")
	}

	chat, err := h.service.GetChat(chatID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get Chat: %v", err)
	}

	messageID, err := h.service.SendMessage(ctx, chatID, senderUsername, receiverUsername, content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	if chat.SourceLang != chat.TargetLang {
		q, err := h.channel.QueueDeclare(broker.MessageSentEvent, true, false, false, false, nil)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to declare queue: %v", err)
		}

		msgData := map[string]interface{}{
			"sourceLang": chat.TargetLang,
			"targetLang": chat.SourceLang,
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
	}

	return &pb.SendMessageResponse{MessageId: messageID}, nil
}

func (h *GrpcHandler) StreamMessages(req *pb.StreamMessagesRequest, stream pb.ChatService_StreamMessagesServer) error {
	chatID := req.GetChatId()

	msgCh, err := h.service.StreamMessages(stream.Context(), chatID)
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
			protoMsg := chatMessageFromModel(msg)
			if err := stream.Send(protoMsg); err != nil {
				return status.Errorf(codes.Internal, "failed to send message: %v", err)
			}
		}
	}
}

func (h *GrpcHandler) GetMessage(_ context.Context, req *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	msg, err := h.service.GetMessage(req.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get message: %v", err)
	}

	return &pb.GetMessageResponse{Message: chatMessageFromModel(msg)}, nil
}

func (h *GrpcHandler) ListMessages(_ context.Context, req *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	if req.GetChatId() == "" {
		return nil, status.Error(codes.InvalidArgument, "chat_id must be provided")
	}

	var since *time.Time
	if req.GetSinceTimestamp() > 0 {
		t := time.Unix(req.GetSinceTimestamp(), 0)
		since = &t
	}

	msgs, pageToken, err := h.service.ListMessages(req.GetChatId(), since, int(req.GetLimit()), req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	protoMsgs := make([]*pb.ChatMessage, len(msgs))
	for i, m := range msgs {
		protoMsgs[i] = chatMessageFromModel(m)
	}

	return &pb.ListMessagesResponse{
		Messages:      protoMsgs,
		NextPageToken: pageToken,
	}, nil
}

func (h *GrpcHandler) GetChat(_ context.Context, req *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	chat, err := h.service.GetChat(req.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get Chat: %v", err)
	}

	return &pb.GetChatResponse{Chat: chatFromModel(chat)}, nil
}

func (h *GrpcHandler) ListChats(_ context.Context, req *pb.ListChatsRequest) (*pb.ListChatsResponse, error) {
	userName := req.GetUsername()
	if userName == "" {
		return nil, status.Error(codes.InvalidArgument, "username must be provided")
	}

	chats, err := h.service.ListChats(userName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list Chats: %v", err)
	}

	protoChats := make([]*pb.Chat, len(chats))
	for i, c := range chats {
		protoChats[i] = chatFromModel(c)
	}

	return &pb.ListChatsResponse{Chats: protoChats}, nil
}
