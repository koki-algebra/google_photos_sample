package main

import "time"

/* ----- MediaItems ----- */

type MediaItem struct {
	ID            string        `json:"id,omitempty"`
	Description   string        `json:"description,omitempty"`
	ProductUrl    string        `json:"productUrl,omitempty"`
	BaseUrl       string        `json:"baseUrl,omitempty"`
	MimeType      string        `json:"mimeType,omitempty"`
	MediaMetadata MediaMetadata `json:"mediaMetadata,omitempty"`
	Filename      string        `json:"filename,omitempty"`
}

type MediaMetadata struct {
	CreationTime *time.Time `json:"creationTime,omitempty"`
	Width        string     `json:"width,omitempty"`
	Height       string     `json:"height,omitempty"`
}

type MediaItems struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

/* ----- Albums -----*/

type Album struct {
	ID                    string `json:"id,omitempty"`
	Title                 string `json:"title,omitempty"`
	ProductUrl            string `json:"productUrl,omitempty"`
	IsWriteable           bool   `json:"isWriteable"`
	MediaItemsCount       string `json:"mediaItemsCount,omitempty"`
	CoverPhotoBaseUrl     string `json:"coverPhotoBaseUrl,omitempty"`
	CoverPhotoMediaItemId string `json:"coverPhotoMediaItemId,omitempty"`
}

type Albums struct {
	Albums        []Album `json:"albums"`
	NextPageToken string  `json:"nextPageToken"`
}

type CreateAlbum struct {
	Album Album `json:"album"`
}
