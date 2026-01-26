package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-demo/recipes-web/internal/repository/memory"
	"github.com/gin-demo/recipes-web/model"
)

var (
	ErrNotFound     = errors.New("recipe not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("recipe conflict")
	ErrPersistence  = errors.New("persistence error")
)

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
	r, err := ctrl.repo.Create(ctx, recipe)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("%w: %v", ErrPersistence, err)
	}

	return r, nil
}

func (ctrl *Controller) GetRecipeByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	recipe, err := ctrl.repo.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, memory.ErrNotFound):
			return model.Recipe{}, ErrNotFound
		default:
			return model.Recipe{}, fmt.Errorf("%w: %v", ErrPersistence, err)
		}
	}

	return recipe, nil
}

func (ctrl *Controller) ListRecipes(ctx context.Context) ([]model.Recipe, error) {
	recipes, err := ctrl.repo.GetAll(ctx)
	if err != nil {
		return nil, ErrPersistence
	}
	return recipes, nil
}

func (ctrl *Controller) UpdateRecipe(ctx context.Context, id model.RecipeID, cmd UpdateRecipeCommand) (model.Recipe, error) {
	existing, err := ctrl.repo.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, memory.ErrNotFound):
			return model.Recipe{}, ErrNotFound
		default:
			return model.Recipe{}, fmt.Errorf("%w: %v", ErrPersistence, err)
		}
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

	r, err := ctrl.repo.Update(ctx, existing)
	if err != nil {
		switch {
		case errors.Is(err, memory.ErrPersistence):
			return model.Recipe{}, ErrPersistence
		default:
			return model.Recipe{}, err
		}
	}

	return r, nil
}

func (ctrl *Controller) DeleteRecipe(ctx context.Context, id model.RecipeID) error {
	err := ctrl.repo.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, memory.ErrPersistence):
			return ErrPersistence
		case errors.Is(err, memory.ErrNotFound):
			return ErrNotFound
		default:
			return ErrPersistence
		}
	}

	return err
}

func (ctrl *Controller) GetRecipeByTag(ctx context.Context, tag string) ([]model.Recipe, error) {
	if tag == "" {
		return []model.Recipe{}, ErrInvalidInput
	}

	return ctrl.repo.GetByTag(ctx, tag)
}
