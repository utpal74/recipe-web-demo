package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for recipe operations.
type Handler struct {
	recipeService domain.RecipeService
}

func New(service domain.RecipeService) *Handler {
	return &Handler{recipeService: service}
}

// CreateRecipeHandler handles POST requests to create a new recipe.
func (handler *Handler) CreateRecipeHandler(ctx *gin.Context) {
	var req createRecipeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	result, err := handler.recipeService.Create(ctx.Request.Context(), toDomainRequest(req))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// ListRecipeHandler handles GET requests to list all recipes.
func (handler *Handler) ListRecipeHandler(ctx *gin.Context) {
	recipes, err := handler.recipeService.GetAll(ctx.Request.Context())
	if err != nil {
		switch {
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("internal error : %v", err)})
		}
		return
	}

	ctx.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler handles PUT requests to update an existing recipe.
func (handler *Handler) UpdateRecipeHandler(ctx *gin.Context) {
	var req updateRecipeIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reciept ID",
		})
		return
	}

	var body updateRecipeRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "need reciept ID",
		})
		return
	}

	updatedRecipe, err := handler.recipeService.Update(ctx.Request.Context(),
		toDomainUpdateRequest(string(req.ID), body))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, updatedRecipe)
}

// ListRecipesByTagHandler handles GET requests to list recipes by tag.
func (handler *Handler) ListRecipesByTagHandler(ctx *gin.Context) {
	var req searchRecipeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "tag is required",
		})
		return
	}

	recipes, err := handler.recipeService.GetByTag(ctx.Request.Context(), req.Tag)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, recipes)
}

// GetRecipeByIDHandler handles GET requests to retrieve a recipe by ID.
func (handler *Handler) GetRecipeByIDHandler(ctx *gin.Context) {
	var req searchByIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reciept ID",
		})
		return
	}

	result, err := handler.recipeService.GetByID(ctx.Request.Context(), req.ID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// DeleteRecipeHandler handles DELETE requests to remove a recipe by ID.
func (handler *Handler) DeleteRecipeHandler(ctx *gin.Context) {
	var req deleteByIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reciept ID",
		})
		return
	}

	if err := handler.recipeService.Delete(ctx.Request.Context(), req.ID); err != nil {
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
