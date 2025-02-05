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
