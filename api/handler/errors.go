package handler

import "errors"

var (
	ErrInvalid  = errors.New("object invalid")
	ErrDup      = errors.New("object duplicated")
	ErrCreate   = errors.New("object creation failed")
	ErrNotFound = errors.New("object not found")
	ErrSave     = errors.New("save config failed")
)
