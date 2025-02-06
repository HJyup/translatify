package main

import (
	"context"
	"github.com/HJyup/translatify-common/discovery"
	"github.com/HJyup/translatify-common/discovery/consul"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/internal/gateway"
	"github.com/HJyup/translatify-gateway/internal/handlers"
	mux2 "github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = "gateway"
	httpAddr    = utils.EnvString("GATEWAY_ADDR")
	consulAddr  = utils.EnvString("CONSUL_ADDR")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err = registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Failed to health check", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.DeRegister(ctx, instanceID, serviceName)

	mux := mux2.NewRouter()

	chatGateway := gateway.NewGateway(registry)
	chatHandler := handlers.NewChatHandler(chatGateway)
	chatHandler.RegisterRoutes(mux)

	log.Println("Starting server on", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
