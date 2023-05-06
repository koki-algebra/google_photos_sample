package photos

import (
	"context"

	"golang.org/x/oauth2"
)

type GooglePhotosService interface {
	CreateAlbum(ctx context.Context, title string) (Album, error)
	GetAlbum(ctx context.Context, albumID string) (Album, error)
	GetAlbums(ctx context.Context, pageToken string) (Albums, error)
	GetAlbumImages(ctx context.Context, albumID string, pageToken string) (MediaItems, error)
	UploadImages(ctx context.Context, mediaItems MediaItems, albumID string) error
	PatchImage(ctx context.Context, mediaItem MediaItem) error
	AddImagesToAlbum(ctx context.Context, albumID string, mediaItems MediaItems) error
}

type googlePhotosServiceImpl struct {
	config        *oauth2.Config
	tokenFilepath string
}

func NewGooglePhotosService(config *oauth2.Config, tokenFilepath string) GooglePhotosService {
	return &googlePhotosServiceImpl{
		config:        config,
		tokenFilepath: tokenFilepath,
	}
}
