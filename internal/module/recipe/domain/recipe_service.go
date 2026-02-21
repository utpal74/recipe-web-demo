package domain

import "context"

type RecipeService interface {
	Create(context.Context, CreateRecipeInput) (Recipe, error)
	GetByID(context.Context, RecipeID) (Recipe, error)
	GetAll(context.Context) ([]Recipe, error)
	Update(context.Context, UpdateRecipeInput) (Recipe, error)
	Delete(context.Context, RecipeID) error
	GetByTag(context.Context, string) ([]Recipe, error)
	CreateMany(ctx context.Context, recipes []Recipe) error
}
