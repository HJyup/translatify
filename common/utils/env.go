package utils

import (
	"log"
	"os"
)

func EnvString(key string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	log.Fatalf("Required environment variable %s is not set or is empty", key)
	return ""
}
