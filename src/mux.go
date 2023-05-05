package main

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

func NewServMux(config *oauth2.Config) (http.Handler, func(), error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Google Photos Library!")
	})
	
	return mux, func() {}, nil
}
