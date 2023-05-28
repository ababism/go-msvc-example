package core

import (
	"errors"
)

type ErrorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

var (
	ErrTokenInvalid  = errors.New("token invalid")
	ErrIncorrectBody = errors.New("incorrect json body")
	ErrInternal      = errors.New("server internal error")
	ErrNotFound      = errors.New("not found")
	ErrAccessDenied  = errors.New("access denied")
	ErrAlreadyExists = errors.New("element already exists")
)

const (
	CodeTokenInvalid  = 1
	CodeIncorrectBody = 2
	CodeInternalError = 3
	CodeNotFound      = 4
	CodeAccessDenied  = 5
	CodeAlreadyExists = 6
)
