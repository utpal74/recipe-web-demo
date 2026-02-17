package domain

import (
	"context"
)

type UserService interface {
	Create(context.Context, CreateUserInput) (User, error)
	Update(context.Context, UpdateUserInput) (User, error)
	Delete(context.Context, UserID) error

	FindByID(context.Context, UserID) (User, error)
	FindByUserName(context.Context, string) (User, error)
}
