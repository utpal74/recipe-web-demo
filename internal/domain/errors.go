package domain

import "errors"

var (
	ErrNotFound      = errors.New("recipe not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrConflict      = errors.New("recipe conflict")
	ErrPersistence   = errors.New("persistence error")
	ErrIOFailure     = errors.New("IO failure")
	ErrSerialization = errors.New("serialization/deserialization failure")
)
