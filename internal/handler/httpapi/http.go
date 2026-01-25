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
	var recipe model.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error creating recipe",
			"details": err.Error(),
		})
		return
	}

	result, err := handler.ctrl.CreateRecipe(ctx.Request.Context(), recipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusCreated, result)
}

func (handler *Handler) ListRecipeHandler(ctx *gin.Context) {
	recipes, err := handler.ctrl.ListRecipes(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error fetching the recipes",
			"details": err.Error(),
		})
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
			"message": "need reciept ID",
			"details": err.Error(),
		})
		return
	}

	var body UpdateRecipeRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "need reciept ID",
			"details": err.Error(),
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, updatedRecipe)
}

type SearchRecipeRequest struct {
	Tag string `form:"tag" binding:"required"`
}

func (handler *Handler) SearchRecipeHandler(ctx *gin.Context) {
	var req SearchRecipeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "tag is required",
			"details": err.Error(),
		})
		return
	}

	recipes, err := handler.ctrl.GetRecipeByTag(ctx.Request.Context(), req.Tag)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error searching using tag",
			"detail":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, recipes)
}

type SearchByIDRequest struct {
	ID model.RecipeID `uri:"id" binding:"required"`
}

func (handler *Handler) SearchRecipeByIDHandler(ctx *gin.Context) {
	var req SearchByIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid reciept ID",
			"details": err.Error(),
		})
		return
	}

	result, err := handler.ctrl.GetRecipeByID(ctx, req.ID)
	if errors.Is(err, recipe.ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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
			"message": "invalid reciept ID",
			"details": err.Error(),
		})
		return
	}

	if err := handler.ctrl.DeleteRecipe(ctx.Request.Context(), req.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error deleting the recipe",
			"error":   err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
