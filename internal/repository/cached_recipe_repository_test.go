package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gin-demo/recipes-web/internal/cache/redisrecipe"
	"github.com/gin-demo/recipes-web/internal/domain"
	"github.com/gin-demo/recipes-web/model"
	"github.com/redis/go-redis/v9"
)

// setupRedisForCachedRepo creates a test Redis client.
func setupRedisForCachedRepo(t *testing.T) (*redis.Client, *redisrecipe.Cache) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping cached repository tests")
	}

	// Clean up before test
	client.FlushDB(ctx)
	cache := redisrecipe.NewCache(client, 1*time.Hour)
	return client, cache
}

// teardownRedisForCachedRepo cleans up the test Redis database.
func teardownRedisForCachedRepo(t *testing.T, client *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.FlushDB(ctx)
	client.Close()
}

// mockRepository implements the recipe repository interface for testing.
type mockRepository struct {
	recipes      []model.Recipe
	createFunc   func(context.Context, model.Recipe) (model.Recipe, error)
	getByIDFunc  func(context.Context, model.RecipeID) (model.Recipe, error)
	getAllFunc   func(context.Context) ([]model.Recipe, error)
	updateFunc   func(context.Context, model.Recipe) (model.Recipe, error)
	deleteFunc   func(context.Context, model.RecipeID) error
	getByTagFunc func(context.Context, string) ([]model.Recipe, error)
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		recipes: make([]model.Recipe, 0),
	}
}

func (m *mockRepository) Create(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, recipe)
	}
	recipe.ID = model.RecipeID("mock-id-" + time.Now().Format("20060102150405"))
	m.recipes = append(m.recipes, recipe)
	return recipe, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	for _, r := range m.recipes {
		if r.ID == id {
			return r, nil
		}
	}
	return model.Recipe{}, domain.ErrNotFound
}

func (m *mockRepository) GetAll(ctx context.Context) ([]model.Recipe, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return m.recipes, nil
}

func (m *mockRepository) Update(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, recipe)
	}
	for i, r := range m.recipes {
		if r.ID == recipe.ID {
			m.recipes[i] = recipe
			return recipe, nil
		}
	}
	return model.Recipe{}, domain.ErrNotFound
}

func (m *mockRepository) Delete(ctx context.Context, id model.RecipeID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	for i, r := range m.recipes {
		if r.ID == id {
			m.recipes = append(m.recipes[:i], m.recipes[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockRepository) GetByTag(ctx context.Context, tag string) ([]model.Recipe, error) {
	if m.getByTagFunc != nil {
		return m.getByTagFunc(ctx, tag)
	}
	return m.recipes, nil
}

func TestNewCachedRepository(t *testing.T) {
	mockRepo := newMockRepository()
	client, cache := setupRedisForCachedRepo(t)
	defer teardownRedisForCachedRepo(t, client)

	cachedRepo := NewCachedRepository(mockRepo, cache)
	if cachedRepo == nil {
		t.Fatal("NewCachedRepository returned nil")
	}
}

func TestCachedRepositoryGetByID_FromCache(t *testing.T) {
	mockRepo := newMockRepository()
	client, cache := setupRedisForCachedRepo(t)
	defer teardownRedisForCachedRepo(t, client)

	recipe := model.Recipe{
		ID:           "recipe-1",
		Name:         "Cached Recipe",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
		PublishedAt:  time.Now(),
	}

	// Pre-populate cache
	_ = cache.SetByID(context.Background(), recipe)

	cachedRepo := NewCachedRepository(mockRepo, cache)

	// GetByID should return from cache without calling repository
	retrieved, err := cachedRepo.GetByID(context.Background(), recipe.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if retrieved.ID != recipe.ID {
		t.Errorf("Expected recipe ID %s, got %s", recipe.ID, retrieved.ID)
	}
	if retrieved.Name != recipe.Name {
		t.Errorf("Expected name %s, got %s", recipe.Name, retrieved.Name)
	}
}

func TestCachedRepositoryGetByID_FromRepository(t *testing.T) {
	mockRepo := newMockRepository()
	client, cache := setupRedisForCachedRepo(t)
	defer teardownRedisForCachedRepo(t, client)

	recipe := model.Recipe{
		ID:           "recipe-2",
		Name:         "Repository Recipe",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
		PublishedAt:  time.Now(),
	}

	// Add recipe to repository
	mockRepo.recipes = append(mockRepo.recipes, recipe)

	cachedRepo := NewCachedRepository(mockRepo, cache)

	// GetByID should retrieve from repository and cache it
	retrieved, err := cachedRepo.GetByID(context.Background(), recipe.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if retrieved.ID != recipe.ID {
		t.Errorf("Expected recipe ID %s, got %s", recipe.ID, retrieved.ID)
	}

	// Verify it was cached
	cachedRecipe, found, _ := cache.GetByID(context.Background(), recipe.ID)
	if !found {
		t.Error("Recipe should be cached after GetByID")
	}
	if cachedRecipe.ID != recipe.ID {
		t.Error("Cached recipe doesn't match retrieved recipe")
	}
}

func TestCachedRepositoryCreate(t *testing.T) {
	mockRepo := newMockRepository()
	client, cache := setupRedisForCachedRepo(t)
	defer teardownRedisForCachedRepo(t, client)

	recipe := model.Recipe{
		Name:         "New Recipe",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
	}

	cachedRepo := NewCachedRepository(mockRepo, cache)

	created, err := cachedRepo.Create(context.Background(), recipe)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.ID == "" {
		t.Error("Created recipe should have an ID")
	}

	// Verify it was cached
	cachedRecipe, found, _ := cache.GetByID(context.Background(), created.ID)
	if !found {
		t.Error("Created recipe should be cached")
	}
	if cachedRecipe.Name != recipe.Name {
		t.Error("Cached recipe doesn't match created recipe")
	}
}

func TestCachedRepositoryUpdate(t *testing.T) {
	mockRepo := newMockRepository()
	client, cache := setupRedisForCachedRepo(t)
	defer teardownRedisForCachedRepo(t, client)

	recipe := model.Recipe{
		ID:           "recipe-3",
		Name:         "Original Name",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
		PublishedAt:  time.Now(),
	}

	// Pre-populate both repository and cache
	mockRepo.recipes = append(mockRepo.recipes, recipe)
	_ = cache.SetByID(context.Background(), recipe)

	cachedRepo := NewCachedRepository(mockRepo, cache)

	// Update recipe
	recipe.Name = "Updated Name"
	updated, err := cachedRepo.Update(context.Background(), recipe)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", updated.Name)
	}

	// Verify cache was invalidated (deleted) for the updated recipe
	_, found, _ := cache.GetByID(context.Background(), recipe.ID)
	if found {
		t.Error("Cache should be invalidated after Update")
	}
}

func TestCachedRepositoryDelete(t *testing.T) {
	mockRepo := newMockRepository()
	client, cache := setupRedisForCachedRepo(t)
	defer teardownRedisForCachedRepo(t, client)

	recipe := model.Recipe{
		ID:           "recipe-4",
		Name:         "Recipe to Delete",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
		PublishedAt:  time.Now(),
	}

	// Pre-populate both repository and cache
	mockRepo.recipes = append(mockRepo.recipes, recipe)
	_ = cache.SetByID(context.Background(), recipe)

	cachedRepo := NewCachedRepository(mockRepo, cache)

	// Verify recipe exists in cache
	_, found, _ := cache.GetByID(context.Background(), recipe.ID)
	if !found {
		t.Fatal("Recipe should exist in cache before deletion")
	}

	// Delete recipe
	err := cachedRepo.Delete(context.Background(), recipe.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it was deleted from repository
	_, err = mockRepo.GetByID(context.Background(), recipe.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Error("Recipe should be deleted from repository")
	}

	// Verify cache was cleared
	_, found, _ = cache.GetByID(context.Background(), recipe.ID)
	if found {
		t.Error("Recipe should be deleted from cache")
	}
}

func TestCachedRepositoryDeleteNotFound(t *testing.T) {
	mockRepo := newMockRepository()
	client, cache := setupRedisForCachedRepo(t)
	defer teardownRedisForCachedRepo(t, client)

	cachedRepo := NewCachedRepository(mockRepo, cache)

	// Try to delete non-existent recipe
	err := cachedRepo.Delete(context.Background(), "non-existent-id")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestCachedRepositoryGetByIDNotFound(t *testing.T) {
	mockRepo := newMockRepository()
	client, cache := setupRedisForCachedRepo(t)
	defer teardownRedisForCachedRepo(t, client)

	cachedRepo := NewCachedRepository(mockRepo, cache)

	// Try to get non-existent recipe
	_, err := cachedRepo.GetByID(context.Background(), "non-existent-id")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
