package domain

import (
	"context"

	"github.com/gin-demo/recipes-web/model"
)

type RecipeRepository interface {
	Create(context.Context, model.Recipe) (model.Recipe, error)
	GetByID(context.Context, model.RecipeID) (model.Recipe, error)
	GetAll(context.Context) ([]model.Recipe, error)
	Update(context.Context, model.Recipe) (model.Recipe, error)
	Delete(context.Context, model.RecipeID) error
	GetByTag(context.Context, string) ([]model.Recipe, error)
}
