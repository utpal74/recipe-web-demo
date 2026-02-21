package domain

import (
	"context"

	"github.com/gin-demo/recipes-web/internal/module/auth/usecase"
)

// AuthService - should responsible to allow that user to sign in, sign out, sign up
type AuthService interface {
	SignIn(ctx context.Context, input usecase.SignInInput) (usecase.SignInOutput, error)
	SignOut(ctx context.Context, identity Identity, refreshToken string) error
	SignUp(ctx context.Context, input usecase.SignUpInput) (usecase.SignUpOutput, error)
	Refresh(ctx context.Context, tokenString string) (string, string, error)
}
