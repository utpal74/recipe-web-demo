package repository

import (
	"context"

	"github.com/gin-demo/recipes-web/internal/cache/redisrecipe"
	"github.com/gin-demo/recipes-web/internal/domain"
	"github.com/gin-demo/recipes-web/model"
)

// CachedRepository wraps a recipe repository with Redis caching layer.
type CachedRepository struct {
	repo  domain.RecipeRepository
	cache *redisrecipe.Cache
}

// NewCachedRepository creates a new CachedRepository with the given repository and cache.
func NewCachedRepository(repo domain.RecipeRepository, cache *redisrecipe.Cache) *CachedRepository {
	return &CachedRepository{repo, cache}
}

// GetByID retrieves a recipe by ID, using the cache when available.
func (c *CachedRepository) GetByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	if r, found, err := c.cache.GetByID(ctx, id); err == nil && found {
		return r, nil
	}

	r, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return model.Recipe{}, err
	}

	_ = c.cache.SetByID(ctx, r)
	return r, nil
}

// Create adds a new recipe and caches the result.
func (c *CachedRepository) Create(ctx context.Context, r model.Recipe) (model.Recipe, error) {
	created, err := c.repo.Create(ctx, r)
	if err != nil {
		return model.Recipe{}, err
	}
	_ = c.cache.SetByID(ctx, created)
	return created, nil
}

// Update modifies a recipe and invalidates its cache entry.
func (c *CachedRepository) Update(ctx context.Context, r model.Recipe) (model.Recipe, error) {
	updated, err := c.repo.Update(ctx, r)
	if err != nil {
		return model.Recipe{}, err
	}
	_ = c.cache.DeleteByID(ctx, r.ID)
	return updated, nil
}

// Delete removes a recipe and clears it from the cache.
func (c *CachedRepository) Delete(ctx context.Context, id model.RecipeID) error {
	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}
	_ = c.cache.DeleteByID(ctx, id)
	return nil
}

// GetAll make a repo call to list down all recipes.
func (c *CachedRepository) GetAll(ctx context.Context) ([]model.Recipe, error) {
	return c.repo.GetAll(ctx)
}

// GetByTag make a repo call to find item based on tag.
func (c *CachedRepository) GetByTag(ctx context.Context, tag string) ([]model.Recipe, error) {
	return c.repo.GetByTag(ctx, tag)
}
