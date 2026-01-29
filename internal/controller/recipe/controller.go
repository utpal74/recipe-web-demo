package recipe

import (
	"context"

	"github.com/gin-demo/recipes-web/internal/domain"
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

// Controller handles business logic for recipe operations.
type Controller struct {
	repo RecipeRepository
}

// New creates a new Controller with the given repository.
func New(repo RecipeRepository) *Controller {
	return &Controller{repo}
}

// CreateRecipe creates a new recipe in the repository.
func (ctrl *Controller) CreateRecipe(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	return ctrl.repo.Create(ctx, recipe)
}

// GetRecipeByID retrieves a recipe by its ID.
func (ctrl *Controller) GetRecipeByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	return ctrl.repo.GetByID(ctx, id)
}

// ListRecipes returns all recipes.
func (ctrl *Controller) ListRecipes(ctx context.Context) ([]model.Recipe, error) {
	return ctrl.repo.GetAll(ctx)
}

// UpdateRecipe updates an existing recipe with the provided command.
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

// DeleteRecipe deletes a recipe by its ID.
func (ctrl *Controller) DeleteRecipe(ctx context.Context, id model.RecipeID) error {
	return ctrl.repo.Delete(ctx, id)
}

// GetRecipeByTag retrieves recipes that have the specified tag.
func (ctrl *Controller) GetRecipeByTag(ctx context.Context, tag string) ([]model.Recipe, error) {
	if tag == "" {
		return []model.Recipe{}, domain.ErrInvalidInput
	}

	return ctrl.repo.GetByTag(ctx, tag)
}
