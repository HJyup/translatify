package main

import (
	common "github.com/HJyup/translatify-common"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
)

var (
	httpAddr = common.EnvString("HTTP_ADDR", ":3000")
)

func main() {
	mux := http.NewServeMux()

	handler := NewHandler()
	handler.registerRoutes(mux)

	log.Printf("Starting server on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start server", err)
	}

}
