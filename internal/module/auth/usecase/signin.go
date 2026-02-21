package usecase

import "time"

type SignInInput struct {
	UserName string
	Password string
}

type SignInOutput struct {
	AccessToken  string
	RefreshToken string
	Expires      time.Time
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

type SignOutInput struct {
	AccessToken string
}
