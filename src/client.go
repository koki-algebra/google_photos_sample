package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type GooglePhotosClient interface {
	CreateAlbum(ctx context.Context, title string) (Album, error)
	GetAlbums(ctx context.Context, pageToken string) (Albums, error)
	GetAlbumImages(ctx context.Context, albumID string, pageToken string) (MediaItems, error)
	UploadImages(ctx context.Context, r io.Reader) error
}

type GooglePhotosClientImpl struct {
	config *oauth2.Config
}

func NewGooglePhotosClient(config *oauth2.Config) GooglePhotosClient {
	return &GooglePhotosClientImpl{
		config: config,
	}
}

func (cli *GooglePhotosClientImpl) CreateAlbum(ctx context.Context, title string) (Album, error) {
	return Album{}, nil
}

func (cli *GooglePhotosClientImpl) GetAlbums(ctx context.Context, pageToken string) (albums Albums, err error) {
	client, err := NewClient(ctx, cli.config, os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}
	token, err := GetTokenFromLocal(os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/albums?pageToken=%s", photosAPIBaseURL, pageToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var myerr MyHTTPError
		if err = json.NewDecoder(resp.Body).Decode(&myerr); err != nil {
			return
		}
		err = NewMyHTTPError(myerr.Code, myerr.Message)
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&albums); err != nil {
		return
	}

	return
}

func (cli *GooglePhotosClientImpl) GetAlbumImages(ctx context.Context, albumID string, pageToken string) (MediaItems, error) {
	return MediaItems{}, nil
}

func (cli *GooglePhotosClientImpl) UploadImages(ctx context.Context, r io.Reader) error {
	return nil
}
