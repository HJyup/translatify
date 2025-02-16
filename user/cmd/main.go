package main

import (
	"context"
	"github.com/HJyup/translatify-user/internal/handler"
	"github.com/HJyup/translatify-user/internal/service"
	"github.com/HJyup/translatify-user/internal/store"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HJyup/translatify-common/tracer"

	"github.com/HJyup/translatify-common/broker"
	"github.com/HJyup/translatify-common/discovery"
	"github.com/HJyup/translatify-common/discovery/consul"
	common "github.com/HJyup/translatify-common/utils"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"

	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = common.EnvString("SERVICE_NAME")
	grpcAddr    = common.EnvString("GRPC_ADDR")
	consulAddr  = common.EnvString("CONSUL_ADDR")

	amqpUser    = common.EnvString("AMQP_USER")
	amqpPass    = common.EnvString("AMQP_PASS")
	amqpHost    = common.EnvString("AMQP_HOST")
	amqpPort    = common.EnvString("AMQP_PORT")
	databaseURL = common.EnvString("DATABASE_URL")

	jaegerAddr = common.EnvString("JAEGER_ADDR")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	if err = registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err = registry.HealthCheck(instanceID, serviceName); err != nil {
					log.Printf("Failed to health check: %v", err)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
	defer registry.DeRegister(ctx, instanceID, serviceName)

	ch, closeConn := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		if err := closeConn(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
		if err := ch.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}()

	config, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	dbConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbConn.Close(ctx)

	grpcServer := grpc.NewServer()
	conn, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
	}
	defer conn.Close()

	str := store.NewStore(dbConn)
	srv := service.NewService(str)
	handler.NewGrpcHandler(grpcServer, srv)

	log.Printf("Starting chat server on %s", grpcAddr)
	if err = grpcServer.Serve(conn); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
