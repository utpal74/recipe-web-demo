package service

import (
	"context"
	"strings"
	"time"

	"github.com/gin-demo/recipes-web/internal/module/auth/service"
	"github.com/gin-demo/recipes-web/internal/module/user/domain"
	"github.com/gin-demo/recipes-web/internal/platform/id"
)

type userService struct {
	userRepo    domain.UserRepository
	pwdHasher   service.PasswordHasher
	idGenerator id.Generator
}

func NewUserService(repo domain.UserRepository,
	pwdHasher service.PasswordHasher,
	idGenerator id.Generator,
) domain.UserService {
	return &userService{
		userRepo:    repo,
		pwdHasher:   pwdHasher,
		idGenerator: idGenerator,
	}
}

// Create implements domain.UserService.
func (u *userService) Create(ctx context.Context, userInput domain.CreateUserInput) (domain.User, error) {
	hash, err := u.pwdHasher.Hash(userInput.Password)
	if err != nil {
		return domain.User{}, domain.ErrInternalServer
	}

	now := time.Now()
	user := domain.NewUser(
		domain.UserID(u.idGenerator.New()),
		strings.TrimSpace(userInput.UserName),
		hash,
		domain.RoleUser,
		now,
	)
	return u.userRepo.Create(ctx, user)
}

// Delete implements domain.UserService.
func (u *userService) Delete(ctx context.Context, id domain.UserID) error {
	return u.userRepo.Delete(ctx, id)
}

// FindByID implements domain.UserService.
func (u *userService) FindByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

// FindByUserName implements domain.UserService.
func (u *userService) FindByUserName(ctx context.Context, userName string) (domain.User, error) {
	return u.userRepo.FindByUserName(ctx, userName)
}

// Update implements domain.UserService.
func (u *userService) Update(ctx context.Context, input domain.UpdateUserInput) (domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, input.ID)
	if err != nil {
		return domain.User{}, domain.ErrNotFound
	}

	if input.UserName != nil {
		user.UserName = *input.UserName
	}

	if input.Password != nil {
		pwdHash, err := u.pwdHasher.Hash(*input.Password)
		if err != nil {
			return domain.User{}, domain.ErrInternalServer
		}
		user.PasswordHash = pwdHash
	}

	user.UpdatedAt = time.Now()

	return u.userRepo.Update(ctx, user)
}
