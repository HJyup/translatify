package handler

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-translation/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedTranslationServiceServer

	service models.TranslationService
	channel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service models.TranslationService, channel *amqp.Channel) {
	handler := &GrpcHandler{
		service: service,
		channel: channel,
	}

	pb.RegisterTranslationServiceServer(grpcServer, handler)
}

func (h *GrpcHandler) TranslateMessage(_ context.Context, req *pb.TranslationRequest) (*pb.TranslationResponse, error) {
	msg, err := h.service.TranslateMessage(req.GetSourceLanguage(), req.GetTargetLanguage(), req.GetContent())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to translate the message: %v", err)
	}

	return &pb.TranslationResponse{
		MessageId:         req.GetMessageId(),
		TranslatedContent: msg.TranslatedContent,
		Success:           true,
	}, nil
}
