package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-demo/recipes-web/internal/auth/domain"
	authdomain "github.com/gin-demo/recipes-web/internal/auth/domain"
	"github.com/gin-demo/recipes-web/internal/auth/usecase"
	userdomain "github.com/gin-demo/recipes-web/internal/user/domain"
)

type authService struct {
	userService  userdomain.UserService
	tokenService TokenService
	pwdHasher    PasswordHasher
}

func NewAuthService(userService userdomain.UserService,
	tokenService TokenService,
	pwdHasher PasswordHasher,
) authdomain.AuthService {
	return &authService{
		userService:  userService,
		tokenService: tokenService,
		pwdHasher:    pwdHasher,
	}
}

// SignIn implements domain.AuthService.
func (a *authService) SignIn(ctx context.Context, request usecase.SignInInput) (usecase.SignInOutput, error) {
	user, err := a.userService.FindByUserName(ctx, request.UserName)
	if err != nil {
		return usecase.SignInOutput{}, domain.ErrInvalidCredentials
	}

	if err := a.pwdHasher.Compare(request.Password, user.PasswordHash); err != nil {
		return usecase.SignInOutput{}, domain.ErrInvalidCredentials
	}

	identity := domain.Identity{
		UserName: user.UserName,
		// Role:     user.Role,
	}

	// TODO: expiry shouldn't be hardcoded but returned by CreateToken here as shown below.
	// token, expiry, err := tokenService.CreateAccessToken(identity)
	expiry := time.Now().Add(15 * time.Minute)

	token, err := a.tokenService.CreateToken(identity, expiry)
	if err != nil {
		return usecase.SignInOutput{}, fmt.Errorf("token creation: %w", err)
	}

	return usecase.SignInOutput{
		Token:   token,
		Expires: expiry,
	}, nil
}

// SignOut implements domain.AuthService.
func (a *authService) SignOut(ctx context.Context, input authdomain.Identity) error {
	panic("unimplemented")
}

// SignUp implements domain.AuthService.
func (a *authService) SignUp(ctx context.Context, input usecase.SignUpInput) (usecase.SignUpOutput, error) {
	in := userdomain.CreateUserInput{
		UserName: input.UserName,
		Password: input.Password,
	}

	user, err := a.userService.Create(ctx, in)
	if err != nil {
		return usecase.SignUpOutput{}, fmt.Errorf("error in user sign up %w", err)
	}

	return usecase.SignUpOutput{
		UserName: user.UserName,
	}, nil
}
