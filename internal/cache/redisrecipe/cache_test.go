package redisrecipe

import (
	"context"
	"testing"
	"time"

	"github.com/gin-demo/recipes-web/model"
	"github.com/redis/go-redis/v9"
)

// setupTestRedis creates a test Redis client pointing to localhost:6379.
func setupTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Clean up before test
	client.FlushDB(ctx)
	return client
}

// teardownTestRedis cleans up the test Redis database.
func teardownTestRedis(t *testing.T, client *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.FlushDB(ctx)
	client.Close()
}

func TestNewCache(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	cache := NewCache(client, 1*time.Hour)
	if cache == nil {
		t.Fatal("NewCache returned nil")
	}
	if cache.client != client {
		t.Error("Cache client not set correctly")
	}
	if cache.ttl != 1*time.Hour {
		t.Error("Cache TTL not set correctly")
	}
}

func TestCacheSetByID(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	cache := NewCache(client, 1*time.Hour)
	ctx := context.Background()

	recipe := model.Recipe{
		ID:           "test-recipe-1",
		Name:         "Test Recipe",
		Tags:         []string{"tag1", "tag2"},
		Ingredients:  []string{"ing1", "ing2"},
		Instructions: []string{"step1", "step2"},
		PublishedAt:  time.Now(),
	}

	err := cache.SetByID(ctx, recipe)
	if err != nil {
		t.Fatalf("SetByID failed: %v", err)
	}

	// Verify the data was stored
	val, err := client.Get(ctx, recipeKey(recipe.ID)).Result()
	if err != nil {
		t.Fatalf("Failed to retrieve cached recipe: %v", err)
	}
	if val == "" {
		t.Error("Recipe not cached properly")
	}
}

func TestCacheGetByID(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	cache := NewCache(client, 1*time.Hour)
	ctx := context.Background()

	recipe := model.Recipe{
		ID:           "test-recipe-2",
		Name:         "Cached Recipe",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
		PublishedAt:  time.Now(),
	}

	// Store recipe in cache
	err := cache.SetByID(ctx, recipe)
	if err != nil {
		t.Fatalf("SetByID failed: %v", err)
	}

	// Retrieve recipe from cache
	retrieved, found, err := cache.GetByID(ctx, recipe.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if !found {
		t.Error("Recipe not found in cache")
	}
	if retrieved.ID != recipe.ID {
		t.Errorf("Expected recipe ID %s, got %s", recipe.ID, retrieved.ID)
	}
	if retrieved.Name != recipe.Name {
		t.Errorf("Expected name %s, got %s", recipe.Name, retrieved.Name)
	}
}

func TestCacheGetByIDNotFound(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	cache := NewCache(client, 1*time.Hour)
	ctx := context.Background()

	// Try to get non-existent recipe
	_, found, _ := cache.GetByID(ctx, "non-existent-id")
	if found {
		t.Error("Expected recipe not to be found")
	}
}

func TestCacheDeleteByID(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	cache := NewCache(client, 1*time.Hour)
	ctx := context.Background()

	recipe := model.Recipe{
		ID:           "test-recipe-3",
		Name:         "Recipe to Delete",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
		PublishedAt:  time.Now(),
	}

	// Store recipe in cache
	err := cache.SetByID(ctx, recipe)
	if err != nil {
		t.Fatalf("SetByID failed: %v", err)
	}

	// Verify it exists
	_, found, _ := cache.GetByID(ctx, recipe.ID)
	if !found {
		t.Fatal("Recipe not found after SetByID")
	}

	// Delete the recipe
	err = cache.DeleteByID(ctx, recipe.ID)
	if err != nil {
		t.Fatalf("DeleteByID failed: %v", err)
	}

	// Verify it's deleted
	_, found, _ = cache.GetByID(ctx, recipe.ID)
	if found {
		t.Error("Recipe should be deleted from cache")
	}
}

func TestCacheTTL(t *testing.T) {
	client := setupTestRedis(t)
	defer teardownTestRedis(t, client)

	cache := NewCache(client, 2*time.Second)
	ctx := context.Background()

	recipe := model.Recipe{
		ID:           "test-recipe-4",
		Name:         "TTL Recipe",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
		PublishedAt:  time.Now(),
	}

	// Store recipe in cache
	err := cache.SetByID(ctx, recipe)
	if err != nil {
		t.Fatalf("SetByID failed: %v", err)
	}

	// Verify it exists
	_, found, _ := cache.GetByID(ctx, recipe.ID)
	if !found {
		t.Fatal("Recipe not found after SetByID")
	}

	// Wait for TTL to expire
	time.Sleep(3 * time.Second)

	// Verify it's expired
	_, found, _ = cache.GetByID(ctx, recipe.ID)
	if found {
		t.Error("Recipe should be expired from cache")
	}
}

func TestRecipeKey(t *testing.T) {
	id := model.RecipeID("test-123")
	key := recipeKey(id)
	expected := "Recipe:test-123"
	if key != expected {
		t.Errorf("Expected key %s, got %s", expected, key)
	}
}
