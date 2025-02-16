package discovery

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
)

func ServiceConnection(serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addr, err := registry.Discover(serviceName)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(
		addr[rand.Intn(len(addr))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
}
