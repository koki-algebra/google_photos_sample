package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
)

type TokenSource struct {
	TokenType    string    `json:"token_type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

func SaveToken(filepath string, token *oauth2.Token) (err error) {
	file, err := os.Create(filepath)
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("defer close error: %v", closeErr)
		}
	}()
	if err != nil {
		return
	}

	data := TokenSource{
		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	if err = json.NewEncoder(file).Encode(&data); err != nil {
		return
	}

	return nil
}

func GetTokenFromLocal(filepath string) (token *oauth2.Token, err error) {
	file, err := os.Open(filepath)
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("defer close error: %v", closeErr)
		}
	}()
	if err != nil {
		return
	}

	var data TokenSource
	if err = json.NewDecoder(file).Decode(&data); err != nil {
		return
	}

	token = &oauth2.Token{
		AccessToken:  data.AccessToken,
		TokenType:    data.TokenType,
		RefreshToken: data.RefreshToken,
		Expiry:       data.Expiry,
	}

	return
}
