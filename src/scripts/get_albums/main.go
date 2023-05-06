package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/koki-algebra/google_photos_sample/auth"
	"github.com/koki-algebra/google_photos_sample/photos"
)

func main() {
	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Panicf("failed to load .env file: %v", err)
	}
	scopes := strings.Split(os.Getenv("SCOPES"), ",")

	// OAuth2 configuration
	config, err := auth.NewGoogleAuthConfig(filepath.Join(os.Getenv("SECRET_DIR"), "client_secret.json"), scopes...)
	if err != nil {
		log.Panicf("failed to initialize google auth config: %v", err)
	}

	// token filepath
	srcTokenFilepath := filepath.Join(os.Getenv("SECRET_DIR"), "src_tokens.json")

	service := photos.NewGooglePhotosService(config, srcTokenFilepath)

	ctx := context.Background()

	var pageToken string
	for {
		albums, err := service.GetAlbums(ctx, pageToken)
		if err != nil {
			log.Panicf("failed to get albums: %v", err)
		}
		pageToken = albums.NextPageToken

		for _, album := range albums.Albums {
			log.Printf("id: %s, title: %s", album.ID, album.Title)
		}

		if pageToken == "" {
			break
		}
	}

	log.Println("finished!")
}
