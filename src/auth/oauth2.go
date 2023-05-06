package auth

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleAuthConfig(filepath string, scope ...string) (*oauth2.Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(data, scope...)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewClient(ctx context.Context, config *oauth2.Config, filepath string) (*http.Client, error) {
	token, err := GetTokenFromLocal(filepath)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, token)

	if err := SaveToken(filepath, token); err != nil {
		return nil, err
	}

	return client, nil
}
