package main

import (
	"log"
)

func main() {
	config, err := NewGoogleAuthConfig("../.secrets/client_secret.json")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("%+v", config)
}
