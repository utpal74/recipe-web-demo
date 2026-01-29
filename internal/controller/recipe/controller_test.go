package recipe

import (
	"context"
	"errors"
	"testing"

	"github.com/gin-demo/recipes-web/internal/repository/memory"
	"github.com/gin-demo/recipes-web/model"
)

type mockRepo struct {
	recipes      []model.Recipe
	createFunc   func(context.Context, model.Recipe) (model.Recipe, error)
	getByIDFunc  func(context.Context, model.RecipeID) (model.Recipe, error)
	getAllFunc   func(context.Context) ([]model.Recipe, error)
	updateFunc   func(context.Context, model.Recipe) (model.Recipe, error)
	deleteFunc   func(context.Context, model.RecipeID) error
	getByTagFunc func(context.Context, string) ([]model.Recipe, error)
}

func (m *mockRepo) Create(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, recipe)
	}
	return recipe, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	for _, r := range m.recipes {
		if r.ID == id {
			return r, nil
		}
	}
	return model.Recipe{}, memory.ErrNotFound
}

func (m *mockRepo) GetAll(ctx context.Context) ([]model.Recipe, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return m.recipes, nil
}

func (m *mockRepo) Update(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, recipe)
	}
	return recipe, nil
}

func (m *mockRepo) Delete(ctx context.Context, id model.RecipeID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockRepo) GetByTag(ctx context.Context, tag string) ([]model.Recipe, error) {
	if m.getByTagFunc != nil {
		return m.getByTagFunc(ctx, tag)
	}
	var result []model.Recipe
	for _, r := range m.recipes {
		for _, t := range r.Tags {
			if t == tag {
				result = append(result, r)
				break
			}
		}
	}
	return result, nil
}

func TestControllerCreateRecipe(t *testing.T) {
	repo := &mockRepo{}
	ctrl := New(repo)

	recipe := model.Recipe{Name: "Test"}
	created, err := ctrl.CreateRecipe(context.Background(), recipe)
	if err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}
	if created.Name != "Test" {
		t.Error("Recipe not created correctly")
	}

	// Error case
	repo.createFunc = func(ctx context.Context, r model.Recipe) (model.Recipe, error) {
		return model.Recipe{}, errors.New("create error")
	}
	_, err = ctrl.CreateRecipe(context.Background(), recipe)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestControllerGetRecipeByID(t *testing.T) {
	repo := &mockRepo{
		recipes: []model.Recipe{{ID: "1", Name: "Test"}},
	}
	ctrl := New(repo)

	recipe, err := ctrl.GetRecipeByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if recipe.ID != "1" {
		t.Error("Wrong recipe")
	}

	// Not found
	_, err = ctrl.GetRecipeByID(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Error("Expected ErrNotFound")
	}

	// Persistence error
	repo.getByIDFunc = func(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
		return model.Recipe{}, errors.New("persistence error")
	}
	_, err = ctrl.GetRecipeByID(context.Background(), "1")
	if err == nil || !errors.Is(err, ErrPersistence) {
		t.Error("Expected ErrPersistence")
	}
}

func TestControllerListRecipes(t *testing.T) {
	repo := &mockRepo{
		recipes: []model.Recipe{{ID: "1"}, {ID: "2"}},
	}
	ctrl := New(repo)

	recipes, err := ctrl.ListRecipes(context.Background())
	if err != nil {
		t.Fatalf("ListRecipes failed: %v", err)
	}
	if len(recipes) != 2 {
		t.Error("Wrong number of recipes")
	}

	// Error case
	repo.getAllFunc = func(ctx context.Context) ([]model.Recipe, error) {
		return nil, errors.New("error")
	}
	_, err = ctrl.ListRecipes(context.Background())
	if err != ErrPersistence {
		t.Error("Expected ErrPersistence")
	}
}

func TestControllerUpdateRecipe(t *testing.T) {
	repo := &mockRepo{
		recipes: []model.Recipe{{ID: "1", Name: "Original"}},
	}
	ctrl := New(repo)

	cmd := UpdateRecipeCommand{
		Name: stringPtr("Updated"),
	}
	updated, err := ctrl.UpdateRecipe(context.Background(), "1", cmd)
	if err != nil {
		t.Fatalf("UpdateRecipe failed: %v", err)
	}
	if updated.Name != "Updated" {
		t.Error("Not updated")
	}

	// Not found
	_, err = ctrl.UpdateRecipe(context.Background(), "nonexistent", cmd)
	if err != ErrNotFound {
		t.Error("Expected ErrNotFound")
	}

	// Persistence error on get
	repo.getByIDFunc = func(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
		return model.Recipe{}, errors.New("get error")
	}
	_, err = ctrl.UpdateRecipe(context.Background(), "1", cmd)
	if err == nil || !errors.Is(err, ErrPersistence) {
		t.Error("Expected ErrPersistence")
	}
}

func TestControllerDeleteRecipe(t *testing.T) {
	repo := &mockRepo{}
	ctrl := New(repo)

	err := ctrl.DeleteRecipe(context.Background(), "1")
	if err != nil {
		t.Fatalf("DeleteRecipe failed: %v", err)
	}

	// Not found
	repo.deleteFunc = func(ctx context.Context, id model.RecipeID) error {
		return memory.ErrNotFound
	}
	err = ctrl.DeleteRecipe(context.Background(), "1")
	if err != ErrNotFound {
		t.Error("Expected ErrNotFound")
	}

	// Persistence error
	repo.deleteFunc = func(ctx context.Context, id model.RecipeID) error {
		return errors.New("persistence")
	}
	err = ctrl.DeleteRecipe(context.Background(), "1")
	if err != ErrPersistence {
		t.Error("Expected ErrPersistence")
	}
}

func TestControllerGetRecipeByTag(t *testing.T) {
	repo := &mockRepo{
		recipes: []model.Recipe{
			{ID: "1", Tags: []string{"a"}},
			{ID: "2", Tags: []string{"b"}},
		},
	}
	ctrl := New(repo)

	recipes, err := ctrl.GetRecipeByTag(context.Background(), "a")
	if err != nil {
		t.Fatalf("GetRecipeByTag failed: %v", err)
	}
	if len(recipes) != 1 {
		t.Error("Wrong number of recipes")
	}

	// Empty tag
	_, err = ctrl.GetRecipeByTag(context.Background(), "")
	if err != ErrInvalidInput {
		t.Error("Expected ErrInvalidInput")
	}
}

func stringPtr(s string) *string {
	return &s
}
