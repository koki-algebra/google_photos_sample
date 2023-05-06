package photos

import (
	"encoding/json"
	"fmt"
	"io"
)

type GooglePhotosError struct {
	Body googlePhotosError `json:"error"`
}

type googlePhotosError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *GooglePhotosError) Error() string {
	return fmt.Sprintf(`{code: %d, message: %s}`, err.Body.Code, err.Body.Message)
}

func NewGooglePhotosError(body io.Reader) error {
	var e GooglePhotosError
	if err := json.NewDecoder(body).Decode(&e); err != nil {
		return err
	}
	return &e
}
