package mongo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gin-demo/recipes-web/internal/module/recipe/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupTestCollection(t *testing.T) (*mongo.Collection, func()) {
	// Skip MongoDB integration tests unless explicitly enabled
	if os.Getenv("ENABLE_MONGO_TESTS") == "" {
		t.Skip("MongoDB integration tests disabled. Set ENABLE_MONGO_TESTS=1 to run")
	}

	if testing.Short() {
		t.Skip("Skipping MongoDB integration tests in short mode")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}

	// Test if we can perform operations
	if err = client.Ping(ctx, nil); err != nil {
		client.Disconnect(context.Background())
		t.Skipf("MongoDB not accessible: %v", err)
	}

	db := client.Database("test_recipes")
	collection := db.Collection("recipes_test_" + t.Name())

	// Drop collection if it exists to avoid validation issues
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	collection.Drop(ctx2)
	cancel2()

	// Cleanup function
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		collection.Drop(ctx)
		client.Disconnect(context.Background())
	}

	return collection, cleanup
}

func TestCreate(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := New(collection)

	recipe := domain.Recipe{
		Name:         "Test Recipe",
		Tags:         []string{"test"},
		Ingredients:  []string{"ing1", "ing2"},
		Instructions: []string{"step1", "step2"},
	}

	created, err := repo.Create(context.Background(), recipe)
	if err != nil {
		// Debug: Print the actual error
		t.Logf("Create returned error: %v (type: %T)", err, err)
		t.Fatalf("Create failed: %v", err)
	}
	if created.ID == "" {
		t.Error("ID not set")
	}
	if created.Name != recipe.Name {
		t.Error("Name not copied")
	}
	if created.PublishedAt.IsZero() {
		t.Error("PublishedAt not set")
	}
}

func TestGetByID(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := New(collection)

	recipe := domain.Recipe{
		Name:         "Test Recipe",
		Tags:         []string{"test"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
	}

	created, _ := repo.Create(context.Background(), recipe)

	// Test GET
	retrieved, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if retrieved.ID != created.ID {
		t.Error("Wrong recipe retrieved")
	}
	if retrieved.Name != recipe.Name {
		t.Error("Name mismatch")
	}

	// Test not found
	_, err = repo.GetByID(context.Background(), domain.RecipeID("nonexistent"))
	if err != domain.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestGetAll(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := New(collection)

	recipe1 := domain.Recipe{Name: "Recipe 1"}
	recipe2 := domain.Recipe{Name: "Recipe 2"}

	repo.Create(context.Background(), recipe1)
	repo.Create(context.Background(), recipe2)

	recipes, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(recipes) < 2 {
		t.Errorf("Expected at least 2 recipes, got %d", len(recipes))
	}
}

func TestUpdate(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := New(collection)

	recipe := domain.Recipe{Name: "Original"}
	created, _ := repo.Create(context.Background(), recipe)

	// Update
	created.Name = "Updated"
	updated, err := repo.Update(context.Background(), created)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != "Updated" {
		t.Error("Name not updated")
	}

	// Verify update
	retrieved, _ := repo.GetByID(context.Background(), created.ID)
	if retrieved.Name != "Updated" {
		t.Error("Update not persisted")
	}
}

func TestDelete(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := New(collection)

	recipe := domain.Recipe{Name: "To Delete"}
	created, _ := repo.Create(context.Background(), recipe)

	// Delete
	err := repo.Delete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(context.Background(), created.ID)
	if err != domain.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestGetByTag(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := New(collection)

	recipe1 := domain.Recipe{Name: "Recipe 1", Tags: []string{"tag1"}}
	recipe2 := domain.Recipe{Name: "Recipe 2", Tags: []string{"tag2"}}

	repo.Create(context.Background(), recipe1)
	repo.Create(context.Background(), recipe2)

	// Get by tag
	recipes, err := repo.GetByTag(context.Background(), "tag1")
	if err != nil {
		t.Fatalf("GetByTag failed: %v", err)
	}
	if len(recipes) == 0 {
		t.Error("No recipes found with tag1")
	}

	found := false
	for _, r := range recipes {
		if r.Name == "Recipe 1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Recipe 1 not found in results")
	}
}

func TestCreateMany(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := New(collection)

	recipes := []domain.Recipe{
		{Name: "Recipe 1", Tags: []string{"test"}},
		{Name: "Recipe 2", Tags: []string{"test"}},
	}

	err := repo.CreateMany(context.Background(), recipes)
	if err != nil {
		t.Fatalf("CreateMany failed: %v", err)
	}

	// Verify they were created
	all, _ := repo.GetAll(context.Background())
	if len(all) < 2 {
		t.Error("CreateMany did not create all recipes")
	}
}
