package httpapi

import (
	"errors"
	"fmt"
	"net/http"

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
