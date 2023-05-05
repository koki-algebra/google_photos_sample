package main

import (
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
