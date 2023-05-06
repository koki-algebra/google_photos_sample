package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type MyHTTPError struct {
	Body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (err *MyHTTPError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", err.Body.Code, err.Body.Message)
}

func NewMyHTTPError(code int, message string) error {
	return &MyHTTPError{
		Body: struct {
			Code    int    "json:\"code\""
			Message string "json:\"message\""
		}{
			Code:    code,
			Message: message,
		},
	}
}

func NewMyHTTPErrorFromReader(body io.Reader) error {
	var myerr MyHTTPError
	if err := json.NewDecoder(body).Decode(&myerr); err != nil {
		return err
	}
	return &myerr
}

func ErrorParser(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var myError *MyHTTPError
	if errors.As(err, &myError) {
		http.Error(w, myError.Body.Message, myError.Body.Code)
		return
	}

	log.Printf("unexpected error: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
