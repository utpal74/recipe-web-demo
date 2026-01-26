package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-demo/recipes-web/internal/controller/recipe"
	"github.com/gin-demo/recipes-web/model"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctrl *recipe.Controller
}

func New(ctrl *recipe.Controller) *Handler {
	return &Handler{ctrl}
}

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

type UpdateRecipeIDRequest struct {
	ID model.RecipeID `uri:"id" binding:"required"`
}

type UpdateRecipeRequest struct {
	Name        *string  `json:"name"`
	Tags        []string `json:"tags"`
	Ingredients []string `json:"ingredients"`
}

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

type SearchRecipeRequest struct {
	Tag string `form:"tag" binding:"required"`
}

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

type SearchByIDRequest struct {
	ID model.RecipeID `uri:"id" binding:"required"`
}

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

type DeleteByIDRequest struct {
	ID model.RecipeID `uri:"id" binding:"required"`
}

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
