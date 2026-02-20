package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-demo/recipes-web/internal/auth/domain"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func New(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (ah *AuthHandler) SignInHandler(ctx *gin.Context) {
	var request signInUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "username or password is required",
		})
		return
	}

	response, err := ah.authService.SignIn(ctx.Request.Context(), RequestToUsecase(request))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrNotAuthorized):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("something went wrong: %v", err),
			})
		}
		return
	}

	ctx.SetSameSite(http.SameSiteStrictMode)

	secure := strings.EqualFold("production", os.Getenv("ENV"))

	ctx.SetCookie(
		"refresh_token",
		response.RefreshToken,
		int(7*24*time.Hour.Seconds()),
		"/auth/refresh",
		"",
		secure,
		true,
	)

	ctx.JSON(http.StatusOK, UsecaseToResponse(response))
}

func (h *AuthHandler) CreateUserHandler(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	reqInput := toCreateUserRequest(req)

	user, err := h.authService.SignUp(ctx.Request.Context(), reqInput)
	if err != nil {
		switch {
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something wrong with user creation",
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, toUserResponse(user.UserName))
}

func (h *AuthHandler) SignOutHandler(ctx *gin.Context) {
	identity := ctx.MustGet("authIdentity").(domain.Identity)

	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	err = h.authService.SignOut(ctx.Request.Context(), identity, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sign out",
		})
		return
	}

	ctx.SetCookie(
		"refresh_token",
		"",
		-1,
		"/auth/refresh",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "signed out"})
}

func (h *AuthHandler) RefreshHandler(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh token required",
		})
		return
	}

	accessToken, refreshToken, err := h.authService.Refresh(ctx.Request.Context(), refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.SetSameSite(http.SameSiteStrictMode)

	secure := strings.EqualFold("production", os.Getenv("ENV"))

	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		int(7*24*time.Hour.Seconds()),
		"/auth/refresh",
		"",
		secure,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}
