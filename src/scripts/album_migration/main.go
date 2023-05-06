package main

import (
	"context"
	"fmt"
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
	dstTokenFilepath := filepath.Join(os.Getenv("SECRET_DIR"), "dst_tokens.json")

	// source service
	srcService := photos.NewGooglePhotosService(config, srcTokenFilepath)
	// destination service
	dstService := photos.NewGooglePhotosService(config, dstTokenFilepath)

	// get album id
	var albumID string
	fmt.Println("Enter target album id...")
	fmt.Scan(&albumID)

	ctx := context.Background()

	// get target album from src
	album, err := srcService.GetAlbum(ctx, albumID)
	if err != nil {
		log.Panicf("failed to get album: %v", err)
	}
	// create album at dst
	if _, err := dstService.CreateAlbum(ctx, album.Title); err != nil {
		log.Panicf("failed to create album: %v", err)
	}

	var pageToken string
	for {
		items, err := srcService.GetAlbumImages(ctx, albumID, pageToken)
		if err != nil {
			log.Panicf("failed to get images: %v", err)
		}
		pageToken = items.NextPageToken

		if pageToken == "" {
			break
		}
	}

	log.Printf("album %s has been migrated successfully", albumID)
}
