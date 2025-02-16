package chat

import (
	"context"
	"log"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/discovery"
)

type GrpcGateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *GrpcGateway {
	return &GrpcGateway{registry: registry}
}

func (g *GrpcGateway) CreateChat(ctx context.Context, payload *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	conn, err := discovery.ServiceConnection("chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.CreateChat(ctx, payload)
}

func (g *GrpcGateway) SendMessage(ctx context.Context, payload *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	conn, err := discovery.ServiceConnection("chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.SendMessage(ctx, payload)
}

func (g *GrpcGateway) GetMessage(ctx context.Context, payload *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	conn, err := discovery.ServiceConnection("chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.GetMessage(ctx, payload)
}

func (g *GrpcGateway) ListMessages(ctx context.Context, payload *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	conn, err := discovery.ServiceConnection("chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.ListMessages(ctx, payload)
}

func (g *GrpcGateway) StreamMessages(ctx context.Context, payload *pb.StreamMessagesRequest) (pb.ChatService_StreamMessagesClient, error) {
	conn, err := discovery.ServiceConnection("chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
		return nil, err
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.StreamMessages(ctx, payload)
}

func (g *GrpcGateway) GetChat(ctx context.Context, payload *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	conn, err := discovery.ServiceConnection("chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.GetChat(ctx, payload)
}

func (g *GrpcGateway) ListChats(ctx context.Context, payload *pb.ListChatsRequest) (*pb.ListChatsResponse, error) {
	conn, err := discovery.ServiceConnection("chat", g.registry)
	if err != nil {
		log.Fatal("Failed to connect to chat service")
	}
	chatClient := pb.NewChatServiceClient(conn)
	return chatClient.ListChats(ctx, payload)
}
