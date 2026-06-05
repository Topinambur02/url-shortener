package exceptions

import "errors"

var (
	ErrNotFound = errors.New("url not found")
	ErrConflict = errors.New("url already exists")
)