package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-demo/recipes-web/model"
)

func TestNew(t *testing.T) {
	// Create temp file with valid JSON
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")

	validRecipes := []model.Recipe{
		{
			ID:           "1",
			Name:         "Recipe 1",
			Tags:         []string{"tag1"},
			Ingredients:  []string{"ing1"},
			Instructions: []string{"step1"},
			PublishedAt:  time.Now(),
		},
	}
	data, _ := json.MarshalIndent(validRecipes, "", " ")
	os.WriteFile(tempFile, data, 0644)

	repo, err := New(tempFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if len(repo.data) != 1 {
		t.Errorf("Expected 1 recipe, got %d", len(repo.data))
	}

	// Test invalid file
	_, err = New("nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}

	// Test invalid JSON
	invalidFile := filepath.Join(tempDir, "invalid.json")
	os.WriteFile(invalidFile, []byte("invalid json"), 0644)
	_, err = New(invalidFile)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestRepositoryCreate(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	os.WriteFile(tempFile, []byte("[]"), 0644)

	repo, _ := New(tempFile)

	recipe := model.Recipe{
		Name:         "New Recipe",
		Tags:         []string{"tag"},
		Ingredients:  []string{"ing"},
		Instructions: []string{"step"},
	}

	created, err := repo.Create(context.Background(), recipe)
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

	// Check persisted
	repo2, _ := New(tempFile)
	if len(repo2.data) != 1 {
		t.Error("Not persisted")
	}
}

func TestRepositoryGetByID(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	recipe := model.Recipe{
		ID:   "test-id",
		Name: "Test",
	}
	data, _ := json.MarshalIndent([]model.Recipe{recipe}, "", " ")
	os.WriteFile(tempFile, data, 0644)

	repo, _ := New(tempFile)

	found, err := repo.GetByID(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found.ID != "test-id" {
		t.Error("Wrong recipe found")
	}

	// Not found
	_, err = repo.GetByID(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Error("Expected ErrNotFound")
	}
}

func TestRepositoryGetAll(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	recipes := []model.Recipe{
		{ID: "1", Name: "R1"},
		{ID: "2", Name: "R2"},
	}
	data, _ := json.MarshalIndent(recipes, "", " ")
	os.WriteFile(tempFile, data, 0644)

	repo, _ := New(tempFile)

	all, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("Expected 2 recipes, got %d", len(all))
	}
	// Check copy, not reference
	all[0].Name = "Modified"
	recheck, _ := repo.GetAll(context.Background())
	if recheck[0].Name == "Modified" {
		t.Error("GetAll returned reference, not copy")
	}
}

func TestRepositoryUpdate(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	recipe := model.Recipe{
		ID:   "1",
		Name: "Original",
	}
	data, _ := json.MarshalIndent([]model.Recipe{recipe}, "", " ")
	os.WriteFile(tempFile, data, 0644)

	repo, _ := New(tempFile)

	updated := model.Recipe{
		ID:   "1",
		Name: "Updated",
	}
	result, err := repo.Update(context.Background(), updated)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if result.Name != "Updated" {
		t.Error("Not updated")
	}

	// Check persisted
	repo2, _ := New(tempFile)
	found, _ := repo2.GetByID(context.Background(), "1")
	if found.Name != "Updated" {
		t.Error("Not persisted")
	}

	// Not found
	_, err = repo.Update(context.Background(), model.Recipe{ID: "nonexistent"})
	if err != ErrNotFound {
		t.Error("Expected ErrNotFound")
	}
}

func TestRepositoryDelete(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	recipes := []model.Recipe{
		{ID: "1", Name: "R1"},
		{ID: "2", Name: "R2"},
	}
	data, _ := json.MarshalIndent(recipes, "", " ")
	os.WriteFile(tempFile, data, 0644)

	repo, _ := New(tempFile)

	err := repo.Delete(context.Background(), "1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Check deleted
	_, err = repo.GetByID(context.Background(), "1")
	if err != ErrNotFound {
		t.Error("Not deleted")
	}

	// Check persisted
	repo2, _ := New(tempFile)
	if len(repo2.data) != 1 {
		t.Error("Not persisted")
	}

	// Not found
	err = repo.Delete(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Error("Expected ErrNotFound")
	}
}

func TestRepositoryGetByTag(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	recipes := []model.Recipe{
		{ID: "1", Name: "R1", Tags: []string{"a", "b"}},
		{ID: "2", Name: "R2", Tags: []string{"b", "c"}},
		{ID: "3", Name: "R3", Tags: []string{"c"}},
	}
	data, _ := json.MarshalIndent(recipes, "", " ")
	os.WriteFile(tempFile, data, 0644)

	repo, _ := New(tempFile)

	found, err := repo.GetByTag(context.Background(), "b")
	if err != nil {
		t.Fatalf("GetByTag failed: %v", err)
	}
	if len(found) != 2 {
		t.Errorf("Expected 2 recipes, got %d", len(found))
	}

	// No matches
	found, _ = repo.GetByTag(context.Background(), "nonexistent")
	if len(found) != 0 {
		t.Error("Expected no recipes")
	}
}

func TestRepositoryConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	os.WriteFile(tempFile, []byte("[]"), 0644)

	repo, _ := New(tempFile)

	ctx := context.Background()

	// Concurrent creates
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			recipe := model.Recipe{Name: fmt.Sprintf("Recipe %d", i)}
			_, err := repo.Create(ctx, recipe)
			if err != nil {
				t.Errorf("Concurrent create failed: %v", err)
			}
			done <- true
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}

	all, _ := repo.GetAll(ctx)
	if len(all) != 10 {
		t.Errorf("Expected 10 recipes, got %d", len(all))
	}
}
