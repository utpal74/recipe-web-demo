package service

import (
	"context"
	"time"

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"github.com/gin-demo/recipes-web/internal/shared/id"
)

// recipeService implements RecipeService for recipe operations.
type recipeService struct {
	idGenerator id.Generator
	repo        domain.RecipeRepository
}

func NewRecipeService(repo domain.RecipeRepository, idGenerator id.Generator) domain.RecipeService {
	// NewRecipeService creates a new recipeService instance.
	return recipeService{
		idGenerator: idGenerator,
		repo:        repo,
	}
}

// Create implements [domain.RecipeService].
func (r recipeService) Create(ctx context.Context, recipe domain.CreateRecipeInput) (domain.Recipe, error) {
	newRecipe := domain.Recipe{
		ID:           domain.RecipeID(r.idGenerator.New()),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  time.Now(),
	}

	return r.repo.Create(ctx, newRecipe)
}

// CreateMany implements [domain.RecipeService].
func (r recipeService) CreateMany(ctx context.Context, recipes []domain.Recipe) error {
	return r.repo.CreateMany(ctx, recipes)
}

// Delete implements [domain.RecipeService].
func (r recipeService) Delete(ctx context.Context, id domain.RecipeID) error {
	return r.repo.Delete(ctx, id)
}

// GetAll implements [domain.RecipeService].
func (r recipeService) GetAll(ctx context.Context) ([]domain.Recipe, error) {
	return r.repo.GetAll(ctx)
}

// GetByID implements [domain.RecipeService].
func (r recipeService) GetByID(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
	return r.repo.GetByID(ctx, id)
}

// GetByTag implements [domain.RecipeService].
func (r recipeService) GetByTag(ctx context.Context, tag string) ([]domain.Recipe, error) {
	if tag == "" {
		return []domain.Recipe{}, domain.ErrInvalidInput
	}

	return r.repo.GetByTag(ctx, tag)
}

// Update implements [domain.RecipeService].
func (r recipeService) Update(ctx context.Context, recipe domain.UpdateRecipeInput) (domain.Recipe, error) {
	existing, err := r.repo.GetByID(ctx, recipe.ID)
	if err != nil {
		return domain.Recipe{}, err
	}

	if recipe.Name != nil {
		existing.Name = *recipe.Name
	}
	if recipe.Tags != nil {
		existing.Tags = recipe.Tags
	}
	if recipe.Ingredients != nil {
		existing.Ingredients = recipe.Ingredients
	}

	return r.repo.Update(ctx, existing)
}
