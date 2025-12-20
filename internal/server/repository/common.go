package repository

import "errors"

// Ошибки уровня репозитория.
var (
	// ErrNoModel модель не найдена.
	ErrNoModel = errors.New("model not found")
	// ErrDataTooBig данные слишком большого размера.
	ErrDataTooBig = errors.New("data too big")
)
