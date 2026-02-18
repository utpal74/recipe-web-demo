package domain

import "errors"

var (
	ErrNotAuthorized      = errors.New("authorization failure")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserLocked         = errors.New("user locked")
	ErrUserAlreadyExists  = errors.New("user exists")
	ErrPasswordWeak       = errors.New("weak password")
	ErrPersistence        = errors.New("refresh token persistence failure")
	ErrNotFound           = errors.New("refresh token not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrHashingToken       = errors.New("token hashing failure")
	ErrExpiredToken       = errors.New("expired token")
	ErrUnAuthorized       = errors.New("user not authorized")
)
