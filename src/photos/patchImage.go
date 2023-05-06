package photos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koki-algebra/google_photos_sample/auth"
)

func (cli *googlePhotosServiceImpl) PatchImage(ctx context.Context, mediaItem MediaItem) error {
	client, err := auth.NewClient(ctx, cli.config, cli.tokenFilepath)
	if err != nil {
		return err
	}
	token, err := auth.GetTokenFromLocal(cli.tokenFilepath)
	if err != nil {
		return err
	}

	var reqBody bytes.Buffer
	if err := json.NewEncoder(&reqBody).Encode(mediaItem); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/mediaItems/%s?updateMask=description", photosAPIBaseURL, mediaItem.ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, &reqBody)
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
