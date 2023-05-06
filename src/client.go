package main

import (
	"bytes"
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

func (cli *GooglePhotosClientImpl) CreateAlbum(ctx context.Context, title string) (album Album, err error) {
	client, err := NewClient(ctx, cli.config, os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}
	token, err := GetTokenFromLocal(os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}

	data := CreateAlbum{Album: Album{Title: title}}
	var reqBody bytes.Buffer
	if err = json.NewEncoder(&reqBody).Encode(data); err != nil {
		return
	}

	url := fmt.Sprintf("%s/albums", photosAPIBaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &reqBody)
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
		err = NewMyHTTPErrorFromReader(resp.Body)
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&album); err != nil {
		return
	}

	return
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
		err = NewMyHTTPErrorFromReader(resp.Body)
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&albums); err != nil {
		return
	}

	return
}

func (cli *GooglePhotosClientImpl) GetAlbumImages(ctx context.Context, albumID string, pageToken string) (items MediaItems, err error) {
	client, err := NewClient(ctx, cli.config, os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}
	token, err := GetTokenFromLocal(os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}

	data := struct {
		AlbumID   string `json:"albumId"`
		PageToken string `json:"pageToken"`
	}{
		AlbumID:   albumID,
		PageToken: pageToken,
	}
	var reqBody bytes.Buffer
	if err = json.NewEncoder(&reqBody).Encode(data); err != nil {
		return
	}

	url := fmt.Sprintf("%s/mediaItems:search", photosAPIBaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &reqBody)
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
		err = NewMyHTTPErrorFromReader(resp.Body)
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return
	}

	return
}

func (cli *GooglePhotosClientImpl) UploadImages(ctx context.Context, r io.Reader) error {
	return nil
}
