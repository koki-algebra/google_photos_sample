package photos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koki-algebra/google_photos_sample/auth"
)

func (cli *googlePhotosServiceImpl) CreateAlbum(ctx context.Context, title string) (album Album, err error) {
	client, err := auth.NewClient(ctx, cli.config, cli.tokenFilepath)
	if err != nil {
		return
	}
	token, err := auth.GetTokenFromLocal(cli.tokenFilepath)
	if err != nil {
		return
	}

	data := CreateAlbum{Album: Album{Title: title}}
	pr, pw := io.Pipe()

	go func() {
		err = json.NewEncoder(pw).Encode(&data)
		pw.Close()
	}()

	url := fmt.Sprintf("%s/albums", photosAPIBaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pr)
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

	if err = json.NewDecoder(resp.Body).Decode(&album); err != nil {
		return
	}

	return
}
