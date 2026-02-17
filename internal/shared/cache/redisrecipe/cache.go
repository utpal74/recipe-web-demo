package redisrecipe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewCache creates a new Cache instance with the given Redis client and TTL.
func NewCache(client *redis.Client, ttl time.Duration) *Cache {
	return &Cache{
		client: client,
		ttl:    ttl,
	}
}

func recipeKey(id domain.RecipeID) string {
	return fmt.Sprintf("Recipe:%s", id)
}

func (c *Cache) GetByID(ctx context.Context, id domain.RecipeID) (domain.Recipe, bool, error) {
	value, err := c.client.Get(ctx, recipeKey(id)).Result()
	if err == redis.Nil {
		return domain.Recipe{}, false, err
	}
	if err != nil {
		return domain.Recipe{}, false, err
	}

	var recipe domain.Recipe
	if err := json.Unmarshal([]byte(value), &recipe); err != nil {
		return domain.Recipe{}, false, err
	}

	return recipe, true, nil
}

func (c *Cache) SetByID(ctx context.Context, recipe domain.Recipe) error {
	data, err := json.Marshal(&recipe)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, recipeKey(recipe.ID), data, c.ttl).Err()
}

// DeleteByID removes a recipe from the cache by ID.
func (c *Cache) DeleteByID(ctx context.Context, id domain.RecipeID) error {
	return c.client.Del(ctx, recipeKey(id)).Err()
}
