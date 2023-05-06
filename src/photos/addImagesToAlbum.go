package photos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koki-algebra/google_photos_sample/auth"
)

func (cli *googlePhotosServiceImpl) AddImagesToAlbum(ctx context.Context, albumID string, mediaItems MediaItems) error {
	client, err := auth.NewClient(ctx, cli.config, cli.tokenFilepath)
	if err != nil {
		return err
	}
	token, err := auth.GetTokenFromLocal(cli.tokenFilepath)
	if err != nil {
		return err
	}

	var (
		data    BatchAddMediaItems
		reqBody bytes.Buffer
	)
	for _, mediaItem := range mediaItems.MediaItems {
		data.MediaItemIds = append(data.MediaItemIds, mediaItem.ID)
	}
	if err := json.NewEncoder(&reqBody).Encode(data); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/albums/%s:batchAddMediaItems", photosAPIBaseURL, albumID)

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
		err = NewGooglePhotosError(resp.Body)
		return err
	}

	return nil
}
