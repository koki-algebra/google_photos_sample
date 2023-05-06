package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/oauth2"
)

func NewServMux(config *oauth2.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	client := NewGooglePhotosClient(config)

	ctrl := NewController(config, client)

	mux.Get("/health", ctrl.Health)
	mux.Get("/auth", ctrl.Auth)
	mux.Get("/callback", ctrl.Callback)

	mux.Post("/albums", ctrl.CreateAlbum)
	mux.Get("/albums", ctrl.GetAlbums)
	mux.Get("/albums/{id}", ctrl.GetAlbumImages)
	mux.Get("/albums/{id}/migration", ctrl.AlbumMigration)

	return mux, func() {}, nil
}
