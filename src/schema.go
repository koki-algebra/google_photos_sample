package main

/* ----- MediaItems ----- */

type MediaItem struct {
	ID            string        `json:"id"`
	Description   string        `json:"description"`
	ProductUrl    string        `json:"productUrl"`
	BaseUrl       string        `json:"baseUrl"`
	MimeType      string        `json:"mimeType"`
	MediaMetadata MediaMetadata `json:"mediaMetadata"`
	Filename      string        `json:"filename"`
}

type MediaMetadata struct {
	CreationTime string `json:"creationTime"`
	Width        string `json:"width"`
	Height       string `json:"height"`
}

type MediaItems struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

/* ----- Albums -----*/

type Album struct {
	ID                    string `json:"id"`
	Title                 string `json:"title"`
	ProductUrl            string `json:"productUrl"`
	IsWriteable           bool   `json:"isWriteable"`
	MediaItemsCount       string `json:"mediaItemsCount"`
	CoverPhotoBaseUrl     string `json:"coverPhotoBaseUrl"`
	CoverPhotoMediaItemId string `json:"coverPhotoMediaItemId"`
}

type Albums struct {
	Albums        []Album `json:"albums"`
	NextPageToken string  `json:"nextPageToken"`
}
