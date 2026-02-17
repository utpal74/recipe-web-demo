package domain

import "errors"

var (
	ErrNotAuthorized      = errors.New("authorization failure")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserLocked         = errors.New("user locked")
	ErrUserAlreadyExists  = errors.New("user exists")
	ErrPasswordWeak       = errors.New("weak password")
)
