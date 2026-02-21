package httpapi

import "time"

type signInUserRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type signInUserResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expiresIn"`
}

type createUserRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Role     string `json:"role"`
}

type signOutResponse struct {
	Message string `json:"message"`
}
