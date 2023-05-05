package main

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

func SaveToken(filepath string, token *oauth2.Token) (err error) {
	file, err := os.Create(filepath)
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("defer close error: %v", closeErr)
		}
	}()
	if err != nil {
		return err
	}

	data := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	if err = json.NewEncoder(file).Encode(&data); err != nil {
		return
	}

	return nil
}
