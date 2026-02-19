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
	config           AuthConfig
}

func NewAuthService(
	userService userdomain.UserService,
	tokenService TokenService,
	pwdHasher PasswordHasher,
	refreshTokenRepo domain.RefreshTokenRepository,
	idGenerator id.Generator,
	config AuthConfig,
) authdomain.AuthService {
	return &authService{
		userService:      userService,
		tokenService:     tokenService,
		pwdHasher:        pwdHasher,
		refreshTokenRepo: refreshTokenRepo,
		idGenerator:      idGenerator,
		config:           config,
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
		UserID: string(user.ID),
		// Role:     user.Role,
	}

	accessExpiry := a.config.AccessTokenTTL
	refreshExpiry := a.config.RefreshTokenTTL

	now := time.Now()

	token, err := a.tokenService.CreateAccessToken(identity, now.Add(accessExpiry))
	if err != nil {
		return usecase.SignInOutput{}, fmt.Errorf("token creation: %w", err)
	}

	fmt.Println("request here ..")
	refreshToken, err := a.tokenService.CreateRefreshToken(identity, now.Add(refreshExpiry))
	if err != nil {
		return usecase.SignInOutput{}, fmt.Errorf("token creation: %w", err)
	}

	newHash := hashToken(refreshToken)

	t := &domain.RefreshToken{
		ID:        authdomain.TokenID(a.idGenerator.New()),
		UserID:    authdomain.UserID(user.ID),
		TokenHash: newHash,
		Revoked:   false,
		CreatedAt: now,
		ExpiresAt: now.Add(refreshExpiry),
	}

	if err := a.refreshTokenRepo.Save(ctx, t); err != nil {
		return usecase.SignInOutput{}, fmt.Errorf("save refresh token: %w", err)
	}

	return usecase.SignInOutput{
		AccessToken:  token,
		RefreshToken: refreshToken,
		Expires:      now.Add(accessExpiry),
	}, nil
}

// SignOut implements domain.AuthService.
func (a *authService) SignOut(ctx context.Context, identity domain.Identity, refreshToken string) error {

	hashed := hashToken(refreshToken)

	token, err := a.refreshTokenRepo.FindByTokenHash(ctx, hashed)
	if err != nil {
		return domain.ErrNotFound
	}

	if token.UserID != authdomain.UserID(identity.UserID) {
		return domain.ErrUnAuthorized
	}

	return a.refreshTokenRepo.DeleteByID(ctx, token.ID)
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

	accessExpiry := a.config.AccessTokenTTL
	refreshExpiry := a.config.RefreshTokenTTL

	now := time.Now()
	accessToken, err := a.tokenService.CreateAccessToken(identity, now.Add(accessExpiry))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.tokenService.CreateRefreshToken(identity, now.Add(refreshExpiry))
	if err != nil {
		return "", "", err
	}

	newHash := hashToken(refreshToken)

	t := &domain.RefreshToken{
		ID:        authdomain.TokenID(a.idGenerator.New()),
		UserID:    storedToken.UserID,
		TokenHash: newHash,
		Revoked:   false,
		CreatedAt: now,
		ExpiresAt: now.Add(refreshExpiry),
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
