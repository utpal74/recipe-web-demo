package domain

import "errors"

var (
	ErrPersistence       = errors.New("persistence failure")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWeakPassword      = errors.New("weak password")
	ErrInternalServer    = errors.New("internal server error")
	ErrNotFound          = errors.New("user not found")
	ErrNothingToUpdate   = errors.New("nothing to update")
)
