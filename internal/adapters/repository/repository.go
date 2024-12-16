package repository

import "errors"

var (
	ErrUploadNotFound = errors.New("upload not found")
	ErrUploadExists   = errors.New("upload already exists")
)
