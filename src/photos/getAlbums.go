package photos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koki-algebra/google_photos_sample/auth"
)

func (cli *googlePhotosServiceImpl) GetAlbums(ctx context.Context, pageToken string) (albums Albums, err error) {
	client, err := auth.NewClient(ctx, cli.config, cli.tokenFilepath)
	if err != nil {
		return
	}
	token, err := auth.GetTokenFromLocal(cli.tokenFilepath)
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
		err = NewGooglePhotosError(resp.Body)
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&albums); err != nil {
		return
	}

	return
}
