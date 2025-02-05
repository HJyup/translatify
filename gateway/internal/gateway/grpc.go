package gateway

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/discovery"
	"log"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) AddMessage(ctx context.Context, payload *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}

	chatClient := pb.NewChatServiceClient(conn)

	return chatClient.SendMessage(ctx, payload)
}

func (g *Gateway) GetMessage(ctx context.Context, payload *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}

	chatClient := pb.NewChatServiceClient(conn)

	return chatClient.GetMessage(ctx, &pb.GetMessageRequest{
		MessageId: payload.MessageId,
	})
}

func (g *Gateway) ListMessages(ctx context.Context, payload *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}

	chatClient := pb.NewChatServiceClient(conn)

	return chatClient.ListMessages(ctx, &pb.ListMessagesRequest{
		UserId:              payload.UserId,
		CorrespondentUserId: payload.CorrespondentUserId,
		SinceTimestamp:      payload.SinceTimestamp,
	})
}

func (g *Gateway) StreamMessages(ctx context.Context, payload *pb.StreamMessagesRequest) (pb.ChatService_StreamMessagesClient, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
		return nil, err
	}

	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.StreamMessages(ctx, payload)
}
