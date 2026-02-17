package domain

import (
	"time"
)

type UserID string

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           UserID
	UserName     string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UpdateUserInput struct {
	ID       UserID
	UserName *string
	Password *string
}

type CreateUserInput struct {
	UserName string
	Password string
}

type SignUpOutput struct {
	UserName string
}

func NewUser(
	id UserID,
	username string,
	passwordHash string,
	role Role,
	now time.Time,
) User {
	return User{
		ID:           id,
		UserName:     username,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
