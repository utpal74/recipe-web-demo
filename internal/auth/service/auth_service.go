package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-demo/recipes-web/internal/auth/domain"
	authdomain "github.com/gin-demo/recipes-web/internal/auth/domain"
	"github.com/gin-demo/recipes-web/internal/auth/usecase"
	"github.com/gin-demo/recipes-web/internal/shared/id"
	userdomain "github.com/gin-demo/recipes-web/internal/user/domain"
)

type authService struct {
	userService      userdomain.UserService
	tokenService     TokenService
	pwdHasher        PasswordHasher
	refreshTokenRepo domain.RefreshTokenRepository
	idGenerator      id.Generator
}

func NewAuthService(
	userService userdomain.UserService,
	tokenService TokenService,
	pwdHasher PasswordHasher,
	refreshTokenRepo domain.RefreshTokenRepository,
	idGenerator id.Generator,
) authdomain.AuthService {
	return &authService{
		userService:      userService,
		tokenService:     tokenService,
		pwdHasher:        pwdHasher,
		refreshTokenRepo: refreshTokenRepo,
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

	token, err := a.tokenService.CreateAccessToken(identity, expiry)
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

// Refresh implements [domain.AuthService].
func (a *authService) Refresh(ctx context.Context, tokenString string) (string, string, error) {
	identity, err := a.tokenService.ValidateRefreshToken(tokenString)
	if err != nil {
		return "", "", domain.ErrInvalidToken
	}

	hashed := hashToken(tokenString)
	if err != nil {
		return "", "", domain.ErrHashingToken
	}

	storedToken, err := a.refreshTokenRepo.FindByTokenHash(ctx, hashed)
	if err != nil {
		return "", "", domain.ErrNotFound
	}

	if storedToken.Revoked || storedToken.ExpiresAt.Before(time.Now()) {
		return "", "", domain.ErrExpiredToken
	}

	err = a.refreshTokenRepo.DeleteByID(ctx, storedToken.ID)
	if err != nil {
		return "", "", err
	}

	accessExpiry := time.Now().Add(15 * time.Minute)
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

	accessToken, err := a.tokenService.CreateAccessToken(identity, accessExpiry)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.tokenService.CreateRefreshToken(identity, refreshExpiry)
	if err != nil {
		return "", "", err
	}

	newHash := hashToken(refreshToken)

	t := &domain.RefreshToken{
		ID:        authdomain.TokenID(a.idGenerator.New()),
		UserID:    storedToken.UserID,
		TokenHash: newHash,
		Revoked:   false,
		CreatedAt: time.Now(),
		ExpiresAt: refreshExpiry,
	}
	a.refreshTokenRepo.Save(ctx, t)

	return accessToken, refreshToken, nil
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

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
