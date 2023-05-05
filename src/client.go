package main

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type MyTokenSource struct {
	src      oauth2.TokenSource
	filepath string
}

func (s *MyTokenSource) Token() (*oauth2.Token, error) {
	token, err := s.src.Token()
	if err != nil {
		return nil, err
	}

	// save token
	if err := SaveToken(s.filepath, token); err != nil {
		return nil, err
	}

	return token, nil
}

func NewClient(ctx context.Context, config *oauth2.Config, filepath string) (*http.Client, error) {
	// get token from local
	token, err := GetTokenFromLocal(filepath)
	if err != nil {
		return nil, err
	}

	tokenSrc := config.TokenSource(ctx, token)
	mySrc := &MyTokenSource{
		src:      tokenSrc,
		filepath: filepath,
	}
	reuseSrc := oauth2.ReuseTokenSource(token, mySrc)

	// refresh token
	client := oauth2.NewClient(ctx, reuseSrc)

	return client, nil
}
