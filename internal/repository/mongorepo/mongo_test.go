package mongorepo

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/gin-demo/recipes-web/internal/domain"
	"github.com/gin-demo/recipes-web/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func setupTestRepo(t *testing.T) *Repository {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	dbName := "test_recipes"

	repo, err := New(uri, dbName)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Test if we can perform operations (check for authentication)
	ctx := context.Background()
	collection := repo.collection(RECIPE_COLLECTION)
	_, err = collection.EstimatedDocumentCount(ctx)
	if err != nil && strings.Contains(err.Error(), "authentication") {
		t.Skip("MongoDB requires authentication. Set MONGO_URI with credentials, e.g., mongodb://user:pass@localhost:27017")
	}

	// Clean up before test
	collection.DeleteMany(ctx, bson.M{})

	return repo
}

func teardownTestRepo(t *testing.T, repo *Repository) {
	ctx := context.Background()
	err := repo.Close(ctx)
	if err != nil {
		t.Errorf("Failed to close repo: %v", err)
	}
}

func TestNew(t *testing.T) {
	repo := setupTestRepo(t)
	defer teardownTestRepo(t, repo)

	if repo.mongoclient == nil {
		t.Error("Client not set")
	}
	if repo.dbName != "test_recipes" {
		t.Error("DB name not set")
	}
}

func TestRepositoryCreate(t *testing.T) {
	repo := setupTestRepo(t)
	defer teardownTestRepo(t, repo)

	ctx := context.Background()
	recipe := model.Recipe{
		Name:         "Test Recipe",
		Tags:         []string{"test"},
		Ingredients:  []string{"ing1", "ing2"},
		Instructions: []string{"step1", "step2"},
	}

	created, err := repo.Create(ctx, recipe)
	if err != nil {
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

func TestRepositoryGetByID(t *testing.T) {
	repo := setupTestRepo(t)
	defer teardownTestRepo(t, repo)

	ctx := context.Background()
	recipe := model.Recipe{
		Name:         "Test Recipe",
		Tags:         []string{"test"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
	}

	created, _ := repo.Create(ctx, recipe)

	retrieved, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if retrieved.ID != created.ID {
		t.Error("ID mismatch")
	}
	if retrieved.Name != recipe.Name {
		t.Error("Name mismatch")
	}

	// Test not found
	_, err = repo.GetByID(ctx, "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestRepositoryGetAll(t *testing.T) {
	repo := setupTestRepo(t)
	defer teardownTestRepo(t, repo)

	ctx := context.Background()
	recipe1 := model.Recipe{Name: "Recipe 1"}
	recipe2 := model.Recipe{Name: "Recipe 2"}

	repo.Create(ctx, recipe1)
	repo.Create(ctx, recipe2)

	recipes, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(recipes) != 2 {
		t.Errorf("Expected 2 recipes, got %d", len(recipes))
	}
}

func TestRepositoryUpdate(t *testing.T) {
	repo := setupTestRepo(t)
	defer teardownTestRepo(t, repo)

	ctx := context.Background()
	recipe := model.Recipe{
		Name:         "Original",
		Tags:         []string{"old"},
		Ingredients:  []string{"ing"},
		Instructions: []string{"step"},
	}

	created, _ := repo.Create(ctx, recipe)

	updatedRecipe := created
	updatedRecipe.Name = "Updated"
	updatedRecipe.Tags = []string{"new"}

	updated, err := repo.Update(ctx, updatedRecipe)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != "Updated" {
		t.Error("Name not updated")
	}
	if len(updated.Tags) != 1 || updated.Tags[0] != "new" {
		t.Error("Tags not updated")
	}
}

func TestRepositoryDelete(t *testing.T) {
	repo := setupTestRepo(t)
	defer teardownTestRepo(t, repo)

	ctx := context.Background()
	recipe := model.Recipe{Name: "To Delete"}

	created, _ := repo.Create(ctx, recipe)

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Check deleted
	_, err = repo.GetByID(ctx, created.ID)
	if err != domain.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}

	// Test delete nonexistent
	err = repo.Delete(ctx, "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected ErrNotFound for nonexistent, got %v", err)
	}
}

func TestRepositoryGetByTag(t *testing.T) {
	repo := setupTestRepo(t)
	defer teardownTestRepo(t, repo)

	ctx := context.Background()
	recipe1 := model.Recipe{Name: "Recipe 1", Tags: []string{"tag1", "common"}}
	recipe2 := model.Recipe{Name: "Recipe 2", Tags: []string{"tag2", "common"}}

	repo.Create(ctx, recipe1)
	repo.Create(ctx, recipe2)

	recipes, err := repo.GetByTag(ctx, "common")
	if err != nil {
		t.Fatalf("GetByTag failed: %v", err)
	}
	if len(recipes) != 2 {
		t.Errorf("Expected 2 recipes, got %d", len(recipes))
	}

	recipes, err = repo.GetByTag(ctx, "tag1")
	if err != nil {
		t.Fatalf("GetByTag failed: %v", err)
	}
	if len(recipes) != 1 {
		t.Errorf("Expected 1 recipe, got %d", len(recipes))
	}
	if recipes[0].Name != "Recipe 1" {
		t.Error("Wrong recipe returned")
	}
}
