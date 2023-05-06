package photos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/koki-algebra/google_photos_sample/auth"
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

type GooglePhotosServiceImpl struct {
	config *oauth2.Config
}

func NewGooglePhotosService(config *oauth2.Config) GooglePhotosService {
	return &GooglePhotosServiceImpl{
		config: config,
	}
}

func (cli *GooglePhotosServiceImpl) CreateAlbum(ctx context.Context, title string) (album Album, err error) {
	client, err := auth.NewClient(ctx, cli.config, os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}
	token, err := auth.GetTokenFromLocal(os.Getenv("TOKENS_FILEPATH"))
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

	if err = json.NewDecoder(resp.Body).Decode(&album); err != nil {
		return
	}

	return
}

func (cli *GooglePhotosServiceImpl) GetAlbum(ctx context.Context, albumID string) (album Album, err error) {
	return
}

func (cli *GooglePhotosServiceImpl) GetAlbums(ctx context.Context, pageToken string) (albums Albums, err error) {
	client, err := auth.NewClient(ctx, cli.config, os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}
	token, err := auth.GetTokenFromLocal(os.Getenv("TOKENS_FILEPATH"))
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
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&albums); err != nil {
		return
	}

	return
}

func (cli *GooglePhotosServiceImpl) GetAlbumImages(ctx context.Context, albumID string, pageToken string) (items MediaItems, err error) {
	client, err := auth.NewClient(ctx, cli.config, os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		return
	}
	token, err := auth.GetTokenFromLocal(os.Getenv("TOKENS_FILEPATH"))
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
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return
	}

	return
}

func (cli *GooglePhotosServiceImpl) UploadImages(ctx context.Context, r io.Reader) error {
	return nil
}

func (cli *GooglePhotosServiceImpl) PatchImage(ctx context.Context, mediaItem MediaItem) (MediaItem, error) {
	return MediaItem{}, nil
}

func (cli *GooglePhotosServiceImpl) AddImagesToAlbum(ctx context.Context, mediaItems MediaItems) error {
	return nil
}
