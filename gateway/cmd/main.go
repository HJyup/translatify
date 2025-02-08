package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HJyup/translatify-common/discovery"
	"github.com/HJyup/translatify-common/discovery/consul"
	"github.com/HJyup/translatify-common/tracer"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/internal/gateway"
	"github.com/HJyup/translatify-gateway/internal/handlers"
	mux2 "github.com/gorilla/mux"

	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = "gateway"
	httpAddr    = utils.EnvString("GATEWAY_ADDR")
	consulAddr  = utils.EnvString("CONSUL_ADDR")

	jaegerAddr = utils.EnvString("JAEGER_ADDR")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	err := tracer.SetGlobalTracer(ctx, serviceName, jaegerAddr)
	if err != nil {
		log.Fatalf("Failed to set global tracer: %v", err)
	}

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		log.Fatalf("Failed to create registry: %v", err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	if err = registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := registry.HealthCheck(instanceID, serviceName); err != nil {
					log.Printf("Failed to health check: %v", err)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
	defer registry.DeRegister(ctx, instanceID, serviceName)

	mux := mux2.NewRouter()

	chatGateway := gateway.NewGateway(registry)
	chatHandler := handlers.NewChatHandler(chatGateway)
	chatHandler.RegisterRoutes(mux)

	log.Println("Starting server on", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
