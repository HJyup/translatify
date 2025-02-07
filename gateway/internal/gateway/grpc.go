package gateway

import (
	"context"
	"log"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) CreateConversation(ctx context.Context, payload *pb.CreateConversationRequest) (*pb.CreateConversationResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.CreateConversation(ctx, payload)
}

func (g *Gateway) SendMessage(ctx context.Context, payload *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
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
	return chatClient.GetMessage(ctx, payload)
}

func (g *Gateway) ListMessages(ctx context.Context, payload *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.ListMessages(ctx, payload)
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

func (g *Gateway) GetConversation(ctx context.Context, payload *pb.GetConversationRequest) (*pb.GetConversationResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.GetConversation(ctx, payload)
}

func (g *Gateway) ListConversations(ctx context.Context, payload *pb.ListConversationsRequest) (*pb.ListConversationsResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.ListConversations(ctx, payload)
}
