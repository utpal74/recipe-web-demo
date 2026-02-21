package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-demo/recipes-web/internal/module/user/domain"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for user operations.
type Handler struct {
	userService domain.UserService
}

func New(service domain.UserService) *Handler {
	// New creates a new Handler for user operations.
	return &Handler{userService: service}
}

func (h *Handler) UpdateUserHandler(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userService.Update(ctx.Request.Context(), toUpdateUserRequest(req))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPersistence):
			ctx.JSON(http.StatusNotModified, gin.H{
				"error": err.Error(),
			})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something failed while updating the resource",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUserHandler(ctx *gin.Context) {
	var delReq reqUserID
	if err := ctx.ShouldBindUri(&delReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	id := toUserID(delReq)
	if err := h.userService.Delete(ctx.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *Handler) FindUserByIDHandler(ctx *gin.Context) {
	var findReq reqUserID
	if err := ctx.ShouldBindUri(&findReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	user, err := h.userService.FindByID(ctx.Request.Context(), domain.UserID(findReq.ID))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) FindUserByNameHandler(ctx *gin.Context) {
	var findReq reqUserName
	if err := ctx.ShouldBindUri(&findReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}

	user, err := h.userService.FindByUserName(ctx.Request.Context(), findReq.Name)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, user)
}
