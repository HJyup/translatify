package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HJyup/translatify-common/discovery"
	"github.com/HJyup/translatify-common/discovery/consul"
	"github.com/HJyup/translatify-common/tracer"
	"github.com/HJyup/translatify-gateway/internal/gateway/chat"
	"github.com/HJyup/translatify-gateway/internal/gateway/user"
	"github.com/HJyup/translatify-gateway/internal/handlers"
	mux2 "github.com/gorilla/mux"
)

type Specification struct {
	ServiceName string `config:"service_name" required:"true" default:"gateway"`
	Address     string `config:"GATEWAY_ADDR" required:"true"`
	Consul      string `config:"CONSUL_ADDR" required:"true"`
	Environment string `config:"ENVIRONMENT" required:"true"`
	Jaeger      string `config:"JAEGER_ADDR" required:"true"`
}

// @title           Translatify API
// @version         1.0
// @description     Chat application with async translation. This API enables users to create chats, send messages, and perform translations asynchronously.

// @termsOfService  http://translatify.io/terms/

// @contact.name   danyil.butov
// @contact.url    https://github.com/HJyup

// @license.name  MIT License
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
// @description Provide your token with `Bearer <token>` format.

// @securityDefinitions.basic  BasicAuth

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to create logger: " + err.Error())
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	var s Specification
	if err = envconfig.Process("gateway", &s); err != nil {
		logger.Fatal("Failed to process environment variables", zap.Error(err))
	}

	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	tracerCfg := tracer.Config{
		ServiceName:    s.ServiceName,
		ServiceVersion: "1.0.0",
		Environment:    s.Environment,
		ExporterAddr:   s.Jaeger,
		Insecure:       false,
		Timeout:        5 * time.Second,
	}
	tp, err := tracer.InitTracer(ctx, tracerCfg)
	if err != nil {
		logger.Fatal("Failed to set global tracer", zap.Error(err))
	}
	defer func() {
		if err = tracer.ShutdownTracer(ctx, tp); err != nil {
			logger.Error("Failed to shutdown tracer", zap.Error(err))
		}
	}()

	registry, err := consul.NewRegistry(s.Consul)
	if err != nil {
		logger.Fatal("Failed to create registry", zap.Error(err))
	}

	instanceID := discovery.GenerateInstanceID(s.ServiceName)
	if err = registry.Register(instanceID, s.ServiceName, s.Address); err != nil {
		logger.Fatal("Failed to register service", zap.Error(err))
	}
	defer func() {
		if err = registry.DeRegister(instanceID); err != nil {
			logger.Error("Failed to deregister service", zap.Error(err))
		}
	}()

	router := mux2.NewRouter()

	chatGateway := chat.NewGateway(registry)
	chatHandler := handlers.NewChatHandler(chatGateway)
	chatHandler.RegisterRoutes(router)

	userGateway := user.NewGateway(registry)
	userHandler := handlers.NewUserHandler(userGateway)
	userHandler.RegisterRoutes(router)

	refHandler := handlers.NewReferenceHandler()
	refHandler.RegisterRoutes(router)

	logger.Info("Starting server", zap.String("address", s.Address))
	if err = http.ListenAndServe(s.Address, router); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
