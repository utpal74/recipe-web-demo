package httpapi

import (
	"github.com/gin-demo/recipes-web/internal/module/recipe/domain"
)

type createRecipeRequest struct {
	Name         string   `json:"name" binding:"required"`
	Tags         []string `json:"tags" binding:"required,min=1"`
	Ingredients  []string `json:"ingredients" binding:"required,min=1"`
	Instructions []string `json:"instructions" binding:"required,min=1"`
}

type updateRecipeIDRequest struct {
	ID domain.RecipeID `uri:"id" binding:"required"`
}

type updateRecipeRequest struct {
	Name         *string  `json:"name"`
	Tags         []string `json:"tags,omitempty"`
	Ingredients  []string `json:"ingredients,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
}

type searchRecipeRequest struct {
	Tag string `form:"tag" binding:"required"`
}

type searchByIDRequest struct {
	ID domain.RecipeID `uri:"id" binding:"required"`
}

type deleteByIDRequest struct {
	ID domain.RecipeID `uri:"id" binding:"required"`
}
