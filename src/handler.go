package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type Controller struct {
	config *oauth2.Config
}

func NewController(config *oauth2.Config) *Controller {
	return &Controller{
		config: config,
	}
}

func (ctrl *Controller) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func (ctrl *Controller) Auth(w http.ResponseWriter, r *http.Request) {
	url := ctrl.config.AuthCodeURL(os.Getenv("OAUTH2_STATE"), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusFound)
}

func (ctrl *Controller) Callback(w http.ResponseWriter, r *http.Request) {
	// get tokens
	code := r.URL.Query().Get("code")
	token, err := ctrl.config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// save access token & refresh token
	if err := SaveToken(os.Getenv("TOKENS_FILEPATH"), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintln(w, "Token successfully saved to local storage")
}

func (ctrl *Controller) GetImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// client
	client, err := NewClient(ctx, ctrl.config, os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get access token
	token, err := GetTokenFromLocal(os.Getenv("TOKENS_FILEPATH"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/mediaItems", photosAPIBaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	// API call
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
