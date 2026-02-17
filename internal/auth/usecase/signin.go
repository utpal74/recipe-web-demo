package usecase

import "time"

type SignInInput struct {
	UserName string
	Password string
}

type SignInOutput struct {
	Token   string
	Expires time.Time
}

type SignUpInput struct {
	Name     string
	Email    string
	Phone    string
	UserName string
	Password string
}

type SignUpOutput struct {
	UserName string
}
