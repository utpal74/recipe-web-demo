package repository

import (
	"context"

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"github.com/gin-demo/recipes-web/internal/shared/cache/redisrecipe"
)

// CachedRepository wraps a repository and adds caching functionality.
// CachedRepository provides recipe persistence with caching.
type CachedRepository struct {
	repo  domain.RecipeRepository
	cache *redisrecipe.Cache
}

func NewCachedRepository(repo domain.RecipeRepository, cache *redisrecipe.Cache) *CachedRepository {
	// NewCachedRepository creates a new CachedRepository with a cache.
	return &CachedRepository{repo, cache}
}

// GetByID retrieves a recipe by ID, using the cache when available.
func (c *CachedRepository) GetByID(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
	if r, found, err := c.cache.GetByID(ctx, id); err == nil && found {
		return r, nil
	}

	r, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Recipe{}, err
	}

	_ = c.cache.SetByID(ctx, r)
	return r, nil
}

// Create adds a new recipe and caches the result.
func (c *CachedRepository) Create(ctx context.Context, r domain.Recipe) (domain.Recipe, error) {
	created, err := c.repo.Create(ctx, r)
	if err != nil {
		return domain.Recipe{}, err
	}
	_ = c.cache.SetByID(ctx, created)
	return created, nil
}

func (c *CachedRepository) CreateMany(ctx context.Context, recipes []domain.Recipe) error {
	return nil
}

// Update modifies a recipe and invalidates its cache entry.
func (c *CachedRepository) Update(ctx context.Context, r domain.Recipe) (domain.Recipe, error) {
	updated, err := c.repo.Update(ctx, r)
	if err != nil {
		return domain.Recipe{}, err
	}
	_ = c.cache.DeleteByID(ctx, r.ID)
	return updated, nil
}

// Delete removes a recipe and clears it from the cache.
func (c *CachedRepository) Delete(ctx context.Context, id domain.RecipeID) error {
	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}
	_ = c.cache.DeleteByID(ctx, id)
	return nil
}

// GetAll make a repo call to list down all recipes.
func (c *CachedRepository) GetAll(ctx context.Context) ([]domain.Recipe, error) {
	return c.repo.GetAll(ctx)
}

// GetByTag make a repo call to find item based on tag.
func (c *CachedRepository) GetByTag(ctx context.Context, tag string) ([]domain.Recipe, error) {
	return c.repo.GetByTag(ctx, tag)
}
