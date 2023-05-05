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

	ctrl := NewController(config)

	mux.HandleFunc("/health", ctrl.Health)
	mux.HandleFunc("/auth", ctrl.Redirect)
	mux.HandleFunc("/callback", ctrl.Callback)

	return mux, func() {}, nil
}
