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

// User represents a user entity in the system.
type User struct {
	ID           UserID
	UserName     string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UpdateUserInput contains fields for updating a user.
type UpdateUserInput struct {
	ID       UserID
	UserName *string
	Password *string
}

// CreateUserInput contains fields for creating a new user.
type CreateUserInput struct {
	UserName string
	Password string
}

// SignUpOutput contains the result of a user sign-up operation.
type SignUpOutput struct {
	UserName string
}

// NewUser creates a new User instance.
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
