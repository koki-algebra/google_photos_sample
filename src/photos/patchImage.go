package photos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	url := fmt.Sprintf("%s/mediaItems/%s?updateMask=description", photosAPIBaseURL, mediaItem.ID)

	pr, pw := io.Pipe()
	go func() {
		err = json.NewEncoder(pw).Encode(&mediaItem)
		pw.Close()
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, pr)
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
