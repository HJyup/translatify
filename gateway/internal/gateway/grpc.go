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

func (g *Gateway) CreateChat(ctx context.Context, payload *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.CreateChat(ctx, payload)
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

func (g *Gateway) GetChat(ctx context.Context, payload *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.GetChat(ctx, payload)
}

func (g *Gateway) ListChats(ctx context.Context, payload *pb.ListChatsRequest) (*pb.ListChatsResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.ListChats(ctx, payload)
}
