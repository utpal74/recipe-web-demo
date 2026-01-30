package recipe

import (
	"context"

	"github.com/gin-demo/recipes-web/internal/domain"
	"github.com/gin-demo/recipes-web/model"
)

// Controller handles business logic for recipe operations.
type Controller struct {
	repo domain.RecipeRepository
}

// New creates a new Controller with the given repository.
func New(repo domain.RecipeRepository) *Controller {
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
