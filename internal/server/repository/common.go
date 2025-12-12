package repository

import "errors"

var (
	ErrNoModel    = errors.New("model not found")
	ErrDataTooBig = errors.New("data too big")
)
