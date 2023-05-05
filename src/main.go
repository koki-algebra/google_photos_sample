package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("failed to terminated server: %v", err)
	}
}

func run(ctx context.Context) error {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load .env file: %v", err)
		return err
	}

	scopes := strings.Split(os.Getenv("SCOPES"), ",")

	config, err := NewGoogleAuthConfig("../.secrets/client_secret.json", scopes...)
	if err != nil {
		log.Printf("failed to initialize google auth config")
		return err
	}

	mux, cleanup, err := NewServMux(config)
	defer cleanup()
	if err != nil {
		log.Printf("failed to initialize router")
		return err
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Printf("server port must be integer")
		return err
	}

	srv := NewServer(mux, port)
	return srv.Run(ctx)
}
