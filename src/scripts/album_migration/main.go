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

const (
	pageSize = 50
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
		mediaItems, err := srcService.GetAlbumImages(ctx, albumID, pageSize, pageToken)
		if err != nil {
			log.Panicf("failed to get images: %v", err)
		}
		pageToken = mediaItems.NextPageToken

		if len(mediaItems.MediaItems) > pageSize {
			log.Panicf("the number of mediaItems must be less than %d", pageSize)
		}

		if err := dstService.UploadImages(ctx, mediaItems, albumID); err != nil {
			log.Panicf("failed to upload images: %v", err)
		}

		// update creationTime
		for _, mediaItem := range mediaItems.MediaItems {
			if err := dstService.PatchImage(ctx, mediaItem); err != nil {
				log.Printf("id = %s error: %v", mediaItem.ID, err)
			}
			log.Printf("id = %s OK", mediaItem.ID)
		}

		if pageToken == "" {
			break
		}
	}

	log.Printf("album %s has been migrated successfully", albumID)
}
