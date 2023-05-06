package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/koki-algebra/google_photos_sample/auth"
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
	if err := auth.SaveToken(os.Getenv("TOKENS_FILEPATH"), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Token successfully saved to local storage")
}
