package discovery

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
)

type Registry interface {
	Register(ctx context.Context, instanceID, serverName, hostPort string) error
	DeRegister(ctx context.Context, instanceID, serverName string) error
	Discover(ctx context.Context, serverName string) ([]string, error)
	HealthCheck(instanceID, serverName string) error
}

func GenerateInstanceID(serverName string) string {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatalf("failed to generate random bytes: %v", err)
	}
	randomHex := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("%s-%s", serverName, randomHex)
}
