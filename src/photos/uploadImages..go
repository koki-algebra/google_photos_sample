package photos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koki-algebra/google_photos_sample/auth"
)

func (cli *googlePhotosServiceImpl) UploadImages(ctx context.Context, mediaItems MediaItems, albumID string) error {
	client, err := auth.NewClient(ctx, cli.config, cli.tokenFilepath)
	if err != nil {
		return err
	}
	token, err := auth.GetTokenFromLocal(cli.tokenFilepath)
	if err != nil {
		return err
	}

	var (
		data    BatchCreateMediaItems
		reqBody bytes.Buffer
	)
	// upload images
	data.AlbumID = albumID
	for _, mediaItem := range mediaItems.MediaItems {
		uploadToken, err := cli.uploadImage(ctx, mediaItem)
		if err != nil {
			return err
		}
		data.NewMediaItems = append(data.NewMediaItems, newMediaItem{
			Description: mediaItem.Description,
			SimpleMediaItem: simpleMediaItem{
				FileName:    mediaItem.Filename,
				UploadToken: uploadToken,
			},
		})
	}
	if err := json.NewEncoder(&reqBody).Encode(data); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/mediaItems:batchCreate", photosAPIBaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewGooglePhotosError(resp.Body)
	}

	return nil
}

func (cli *googlePhotosServiceImpl) uploadImage(ctx context.Context, mediaItem MediaItem) (uploadToken string, err error) {
	client, err := auth.NewClient(ctx, cli.config, cli.tokenFilepath)
	if err != nil {
		return
	}
	token, err := auth.GetTokenFromLocal(cli.tokenFilepath)
	if err != nil {
		return
	}

	// get image
	resp, err := http.Get(mediaItem.BaseUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// image byte data
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/uploads", photosAPIBaseURL)

	// image upload request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Content-type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-Content-Type", mediaItem.MimeType)
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	uploadResp, err := client.Do(req)
	if err != nil {
		return
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK {
		err = NewGooglePhotosError(uploadResp.Body)
		return
	}

	buf = new(bytes.Buffer)
	if _, err = buf.ReadFrom(uploadResp.Body); err != nil {
		return
	}
	uploadToken = buf.String()

	return
}
