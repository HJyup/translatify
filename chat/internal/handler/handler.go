package handler

import (
	"context"
	"github.com/HJyup/translatify-chat/internal/models"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
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

func (h *GrpcHandler) SendMessage(context.Context, *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}

func (h *GrpcHandler) StreamMessages(*pb.StreamMessagesRequest, grpc.ServerStreamingServer[pb.ChatMessage]) error {
	return status.Errorf(codes.Unimplemented, "method StreamMessages not implemented")
}

func (h *GrpcHandler) GetMessage(ctx context.Context, p *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	log.Println("Are u trying to get a message?")
	msg := &pb.ChatMessage{
		MessageId:         "25",
		FromUserId:        "22",
		ToUserId:          "34",
		Content:           "Hello, World!",
		TranslatedContent: "Hallo, Welt!",
		Timestamp:         234234,
		Translated:        false,
	}

	q, err := h.channel.QueueDeclare(broker.MessageSentEvent, true, false, false, false, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to declare queue: %v", err)
	}

	err = h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         []byte("25"),
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	}

	return &pb.GetMessageResponse{Message: msg}, nil
}
func (h *GrpcHandler) ListMessages(context.Context, *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMessages not implemented")
}
