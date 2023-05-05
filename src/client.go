package main

import (
	"context"
	"io"
	"net/http"
)

type GooglePhotosClient interface {
	CreateAlbum(ctx context.Context, title string) (Album, error)
	GetAlbums(ctx context.Context, pageToken string) (Albums, error)
	GetAlbumImages(ctx context.Context, albumID string, pageToken string) (MediaItems, error)
	UploadImages(ctx context.Context, r io.Reader) error
}

type GooglePhotosClientImpl struct {
	client *http.Client
}

func NewGooglePhotosClient(client *http.Client) GooglePhotosClient {
	return &GooglePhotosClientImpl{
		client: client,
	}
}

func (cli *GooglePhotosClientImpl) CreateAlbum(ctx context.Context, title string) (Album, error) {
	return Album{}, nil
}

func (cli *GooglePhotosClientImpl) GetAlbums(ctx context.Context, pageToken string) (Albums, error) {
	return Albums{}, nil
}

func (cli *GooglePhotosClientImpl) GetAlbumImages(ctx context.Context, albumID string, pageToken string) (MediaItems, error) {
	return MediaItems{}, nil
}

func (cli *GooglePhotosClientImpl) UploadImages(ctx context.Context, r io.Reader) error {
	return nil
}
