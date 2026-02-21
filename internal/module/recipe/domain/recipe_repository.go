package domain

import "context"

type RecipeRepository interface {
	Create(context.Context, Recipe) (Recipe, error)
	GetByID(context.Context, RecipeID) (Recipe, error)
	GetAll(context.Context) ([]Recipe, error)
	Update(context.Context, Recipe) (Recipe, error)
	Delete(context.Context, RecipeID) error
	GetByTag(context.Context, string) ([]Recipe, error)
	CreateMany(ctx context.Context, recipes []Recipe) error
}
