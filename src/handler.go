package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

type Controller struct {
	config *oauth2.Config
	client GooglePhotosClient
}

func NewController(config *oauth2.Config, client GooglePhotosClient) *Controller {
	return &Controller{
		config: config,
		client: client,
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
		return
	}

	fmt.Fprintln(w, "Token successfully saved to local storage")
}

func (ctrl *Controller) GetAlbums(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pageToken := r.URL.Query().Get("pageToken")

	albums, err := ctrl.client.GetAlbums(ctx, pageToken)
	if err != nil {
		ErrorParser(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, albums)
}

func (ctrl *Controller) GetAlbumImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	pageToken := r.URL.Query().Get("pageToken")

	imgs, err := ctrl.client.GetAlbumImages(ctx, id, pageToken)
	if err != nil {
		ErrorParser(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, imgs)
}

func (ctrl *Controller) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body CreateAlbum
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ErrorParser(w, err)
		return
	}

	album, err := ctrl.client.CreateAlbum(ctx, body.Album.Title)
	if err != nil {
		ErrorParser(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, album)
}

func (ctrl *Controller) AlbumMigration(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Album Migration")

	/*
		createAlbum(albumID)
		for {
			hasNext, images := getImages(albumID)
			uploadImages(images)
			for each image of images {  // update creationTime
				patchImage(image)
			}
			addImagesToAlbum(images)
			if !hasNext {
				break
			}
		}
	*/
}
