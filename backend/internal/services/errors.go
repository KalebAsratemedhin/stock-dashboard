package services

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal error")
)