package domain

import "context"

// UserRepository defines methods for user persistence operations.
type UserRepository interface {
	Create(context.Context, User) (User, error)
	Update(context.Context, User) (User, error)
	Delete(context.Context, UserID) error

	FindByID(context.Context, UserID) (User, error)
	FindByUserName(context.Context, string) (User, error)
}
