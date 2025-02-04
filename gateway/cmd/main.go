package main

import (
	common "github.com/HJyup/translatify-common"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-gateway/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
)

var (
	httpAddr        = common.EnvString("GATEWAY_ADDR", ":1234")
	chatServiceAddr = common.EnvString("CHAT_SERVICE_ADDR", "localhost:50051")
)

func main() {
	conn, err := grpc.NewClient(chatServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to chat service", err)
	}
	defer conn.Close()

	mux := http.NewServeMux()

	chatClient := pb.NewChatServiceClient(conn)

	chatHandler := handlers.NewChatHandler(chatClient)
	chatHandler.RegisterRoutes(mux)

	log.Println("Starting server on", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
