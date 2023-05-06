package photos

import (
	"context"
	"io"

	"golang.org/x/oauth2"
)

type GooglePhotosService interface {
	CreateAlbum(ctx context.Context, title string) (Album, error)
	GetAlbum(ctx context.Context, albumID string) (Album, error)
	GetAlbums(ctx context.Context, pageToken string) (Albums, error)
	GetAlbumImages(ctx context.Context, albumID string, pageToken string) (MediaItems, error)
	UploadImages(ctx context.Context, r io.Reader) error
	PatchImage(ctx context.Context, mediaItem MediaItem) (MediaItem, error)
	AddImagesToAlbum(ctx context.Context, mediaItems MediaItems) error
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

func (cli *googlePhotosServiceImpl) UploadImages(ctx context.Context, r io.Reader) error {
	return nil
}

func (cli *googlePhotosServiceImpl) PatchImage(ctx context.Context, mediaItem MediaItem) (MediaItem, error) {
	return MediaItem{}, nil
}

func (cli *googlePhotosServiceImpl) AddImagesToAlbum(ctx context.Context, mediaItems MediaItems) error {
	return nil
}
