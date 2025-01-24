package main

import (
	common "github.com/HJyup/translatify-common"
	pb "github.com/HJyup/translatify-common/api"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

var (
	httpAddr        = common.EnvString("HTTP_ADDR", ":3000")
	userServiceAddr = common.EnvString("USER_SERVICE_ADDR", "localhost:50051")
)

func main() {
	conn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to user service", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	mux := http.NewServeMux()

	handler := NewHandler(client)
	handler.registerRoutes(mux)

	log.Printf("Starting server on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
