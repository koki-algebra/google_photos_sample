package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type MyHTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *MyHTTPError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", err.Code, err.Message)
}

func NewMyHTTPError(code int, message string) error {
	return &MyHTTPError{
		Code:    code,
		Message: message,
	}
}

func ErrorParser(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var myError *MyHTTPError
	if errors.As(err, &myError) {
		http.Error(w, myError.Message, myError.Code)
		return
	}

	log.Printf("unexpected error: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
