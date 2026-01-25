package recipe

import (
	"context"
	"errors"

	"github.com/gin-demo/recipes-web/model"
)

var ErrNotFound = errors.New("recipe not found")

type recipeRepository interface {
	Create(context.Context, model.Recipe) (model.Recipe, error)
	GetByID(context.Context, model.RecipeID) (model.Recipe, error)
	GetAll(context.Context) ([]model.Recipe, error)
	Update(context.Context, model.Recipe) (model.Recipe, error)
	Delete(context.Context, model.RecipeID) error
	GetByTag(context.Context, string) ([]model.Recipe, error)
}

type Controller struct {
	repo recipeRepository
}

func New(repo recipeRepository) *Controller {
	return &Controller{repo}
}

func (ctrl *Controller) CreateRecipe(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	return ctrl.repo.Create(ctx, recipe)
}

func (ctrl *Controller) GetRecipeByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	return ctrl.repo.GetByID(ctx, id)
}

func (ctrl *Controller) ListRecipes(ctx context.Context) ([]model.Recipe, error) {
	return ctrl.repo.GetAll(ctx)
}

func (ctrl *Controller) UpdateRecipe(ctx context.Context, id model.RecipeID, cmd UpdateRecipeCommand) (model.Recipe, error) {
	existing, err := ctrl.repo.GetByID(ctx, id)
	if err != nil {
		return model.Recipe{}, err
	}

	if cmd.Name != nil {
		existing.Name = *cmd.Name
	}
	if cmd.Tags != nil {
		existing.Tags = cmd.Tags
	}
	if cmd.Ingredients != nil {
		existing.Ingredients = cmd.Ingredients
	}

	return ctrl.repo.Update(ctx, existing)
}

func (ctrl *Controller) DeleteRecipe(ctx context.Context, id model.RecipeID) error {
	return ctrl.repo.Delete(ctx, id)
}

func (ctrl *Controller) GetRecipeByTag(ctx context.Context, tag string) ([]model.Recipe, error) {
	return ctrl.repo.GetByTag(ctx, tag)
}
