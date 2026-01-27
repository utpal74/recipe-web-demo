package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-demo/recipes-web/internal/controller/recipe"
	"github.com/gin-demo/recipes-web/model"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for recipe operations.
type Handler struct {
	ctrl *recipe.Controller
}

// New creates a new Handler with the given controller.
func New(ctrl *recipe.Controller) *Handler {
	return &Handler{ctrl}
}

// CreateRecipeHandler handles POST requests to create a new recipe.
func (handler *Handler) CreateRecipeHandler(ctx *gin.Context) {
	var r model.Recipe
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	result, err := handler.ctrl.CreateRecipe(ctx.Request.Context(), r)
	if err != nil {
		switch {
		case errors.Is(err, recipe.ErrInvalidInput):
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
	recipes, err := handler.ctrl.ListRecipes(ctx)
	if err != nil {
		switch {
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, recipes)
}

// UpdateRecipeIDRequest represents the URI parameters for updating a recipe.
type UpdateRecipeIDRequest struct {
	// ID is the unique identifier of the recipe to update
	ID model.RecipeID `uri:"id" binding:"required"`
}

// UpdateRecipeRequest represents the request body for updating a recipe.
type UpdateRecipeRequest struct {
	// Name is the optional new name for the recipe
	Name        *string  `json:"name"`
	// Tags is the optional new list of tags for the recipe
	Tags        []string `json:"tags"`
	// Ingredients is the optional new list of ingredients for the recipe
	Ingredients []string `json:"ingredients"`
}

// UpdateRecipeHandler handles PUT requests to update an existing recipe.
func (handler *Handler) UpdateRecipeHandler(ctx *gin.Context) {
	var req UpdateRecipeIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reciept ID",
		})
		return
	}

	var body UpdateRecipeRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "need reciept ID",
		})
		return
	}

	cmd := recipe.UpdateRecipeCommand{
		Name:        body.Name,
		Tags:        body.Tags,
		Ingredients: body.Ingredients,
	}

	updatedRecipe, err := handler.ctrl.UpdateRecipe(ctx.Request.Context(), req.ID, cmd)
	if err != nil {
		switch {
		case errors.Is(err, recipe.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, updatedRecipe)
}

// SearchRecipeRequest represents the query parameters for searching recipes by tag.
type SearchRecipeRequest struct {
	// Tag is the tag to search recipes by
	Tag string `form:"tag" binding:"required"`
}

// ListRecipesByTagHandler handles GET requests to list recipes by tag.
func (handler *Handler) ListRecipesByTagHandler(ctx *gin.Context) {
	var req SearchRecipeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "tag is required",
		})
		return
	}

	recipes, err := handler.ctrl.GetRecipeByTag(ctx.Request.Context(), req.Tag)
	if err != nil {
		switch {
		case errors.Is(err, recipe.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, recipes)
}

// SearchByIDRequest represents the URI parameters for searching a recipe by ID.
type SearchByIDRequest struct {
	// ID is the unique identifier of the recipe to retrieve
	ID model.RecipeID `uri:"id" binding:"required"`
}

// GetRecipeByIDHandler handles GET requests to retrieve a recipe by ID.
func (handler *Handler) GetRecipeByIDHandler(ctx *gin.Context) {
	var req SearchByIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reciept ID",
		})
		return
	}

	result, err := handler.ctrl.GetRecipeByID(ctx.Request.Context(), req.ID)
	if err != nil {
		switch {
		case errors.Is(err, recipe.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, recipe.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// DeleteByIDRequest represents the URI parameters for deleting a recipe by ID.
type DeleteByIDRequest struct {
	// ID is the unique identifier of the recipe to delete
	ID model.RecipeID `uri:"id" binding:"required"`
}

// DeleteRecipeHandler handles DELETE requests to remove a recipe by ID.
func (handler *Handler) DeleteRecipeHandler(ctx *gin.Context) {
	var req DeleteByIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reciept ID",
		})
		return
	}

	if err := handler.ctrl.DeleteRecipe(ctx.Request.Context(), req.ID); err != nil {
		switch {
		case errors.Is(err, recipe.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
