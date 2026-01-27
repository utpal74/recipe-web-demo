package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-demo/recipes-web/internal/controller/recipe"
	"github.com/gin-demo/recipes-web/internal/repository/memory"
	"github.com/gin-demo/recipes-web/model"
	"github.com/gin-gonic/gin"
)

type mockRepo struct {
	createFunc  func(context.Context, model.Recipe) (model.Recipe, error)
	getByIDFunc func(context.Context, model.RecipeID) (model.Recipe, error)
	listFunc    func(context.Context) ([]model.Recipe, error)
	updateFunc  func(context.Context, model.Recipe) (model.Recipe, error)
	deleteFunc  func(context.Context, model.RecipeID) error
	getByTagFunc func(context.Context, string) ([]model.Recipe, error)
}

func (m *mockRepo) Create(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, recipe)
	}
	recipe.ID = "test-id"
	return recipe, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return model.Recipe{ID: id, Name: "Test"}, nil
}

func (m *mockRepo) GetAll(ctx context.Context) ([]model.Recipe, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return []model.Recipe{}, nil
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
	return []model.Recipe{}, nil
}

// Helper to setup router with handlers
func setupTestRouter(repo *mockRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	ctrl := recipe.New(repo)
	handler := New(ctrl)

	router.GET("/recipes", handler.ListRecipeHandler)
	router.GET("/recipes/search", handler.ListRecipesByTagHandler)
	router.GET("/recipes/:id", handler.GetRecipeByIDHandler)
	router.POST("/recipes", handler.CreateRecipeHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)

	return router
}

func TestCreateRecipeHandler(t *testing.T) {
	repo := &mockRepo{}
	router := setupTestRouter(repo)

	recipe := model.Recipe{Name: "Test Recipe", Tags: []string{"tag1"}}
	body, _ := json.Marshal(recipe)

	req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Invalid JSON
	req2, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer([]byte("invalid")))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w2.Code)
	}

	// Persistence error
	repo.createFunc = func(ctx context.Context, r model.Recipe) (model.Recipe, error) {
		return model.Recipe{}, errors.New("persistence error")
	}
	req3, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer(body))
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w3.Code)
	}
}

func TestListRecipeHandler(t *testing.T) {
	repo := &mockRepo{
		listFunc: func(ctx context.Context) ([]model.Recipe, error) {
			return []model.Recipe{{ID: "1", Name: "R1"}, {ID: "2", Name: "R2"}}, nil
		},
	}
	router := setupTestRouter(repo)

	req, _ := http.NewRequest("GET", "/recipes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Error case
	repo.listFunc = func(ctx context.Context) ([]model.Recipe, error) {
		return nil, errors.New("error")
	}
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)

	if w2.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w2.Code)
	}
}

func TestUpdateRecipeHandler(t *testing.T) {
	repo := &mockRepo{}
	router := setupTestRouter(repo)

	cmd := UpdateRecipeRequest{Name: stringPtr("Updated")}
	body, _ := json.Marshal(cmd)

	req, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Invalid JSON body
	req2, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer([]byte("invalid")))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w2.Code)
	}

	// Not found
	repo.updateFunc = func(ctx context.Context, r model.Recipe) (model.Recipe, error) {
		return model.Recipe{}, recipe.ErrNotFound
	}
	req3, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer(body))
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w3.Code)
	}

	// Persistence error
	repo.updateFunc = func(ctx context.Context, r model.Recipe) (model.Recipe, error) {
		return model.Recipe{}, recipe.ErrPersistence
	}
	req4, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer(body))
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w4.Code)
	}
}

func TestListRecipesByTagHandler(t *testing.T) {
	repo := &mockRepo{
		getByTagFunc: func(ctx context.Context, tag string) ([]model.Recipe, error) {
			if tag == "test" {
				return []model.Recipe{{ID: "1", Tags: []string{"test"}}}, nil
			}
			return []model.Recipe{}, nil
		},
	}
	router := setupTestRouter(repo)

	req, _ := http.NewRequest("GET", "/recipes/search?tag=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Missing tag
	req2, _ := http.NewRequest("GET", "/recipes/search", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w2.Code)
	}

	// Invalid input
	repo.getByTagFunc = func(ctx context.Context, tag string) ([]model.Recipe, error) {
		return nil, recipe.ErrInvalidInput
	}
	req3, _ := http.NewRequest("GET", "/recipes/search?tag=empty", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w3.Code)
	}
}

func TestGetRecipeByIDHandler(t *testing.T) {
	repo := &mockRepo{
		getByIDFunc: func(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
			if id == "1" {
				return model.Recipe{ID: "1", Name: "Test"}, nil
			}
			return model.Recipe{}, memory.ErrNotFound
		},
	}
	router := setupTestRouter(repo)

	req, _ := http.NewRequest("GET", "/recipes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Not found - but controller converts it to ErrNotFound
	repo.getByIDFunc = func(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
		return model.Recipe{}, memory.ErrNotFound
	}
	req2, _ := http.NewRequest("GET", "/recipes/nonexistent", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w2.Code)
	}

	// Persistence error
	repo.getByIDFunc = func(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
		return model.Recipe{}, memory.ErrPersistence
	}
	req3, _ := http.NewRequest("GET", "/recipes/1", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w3.Code)
	}
}

func TestDeleteRecipeHandler(t *testing.T) {
	repo := &mockRepo{
		deleteFunc: func(ctx context.Context, id model.RecipeID) error {
			if id == "1" {
				return nil
			}
			return memory.ErrNotFound
		},
	}
	router := setupTestRouter(repo)

	req, _ := http.NewRequest("DELETE", "/recipes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Not found
	repo.deleteFunc = func(ctx context.Context, id model.RecipeID) error {
		return memory.ErrNotFound
	}
	req2, _ := http.NewRequest("DELETE", "/recipes/nonexistent", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w2.Code)
	}

	// Persistence error
	repo.deleteFunc = func(ctx context.Context, id model.RecipeID) error {
		return memory.ErrPersistence
	}
	req3, _ := http.NewRequest("DELETE", "/recipes/1", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w3.Code)
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	repo := &mockRepo{}
	router := setupTestRouter(repo)

	// Empty recipe name
	recipe := model.Recipe{Name: "", Tags: []string{}}
	body, _ := json.Marshal(recipe)
	req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Should still create even with empty name
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 for empty name, got %d", w.Code)
	}

	// Large lists
	repo.listFunc = func(ctx context.Context) ([]model.Recipe, error) {
		recipes := make([]model.Recipe, 1000)
		for i := 0; i < 1000; i++ {
			recipes[i] = model.Recipe{ID: model.RecipeID("id" + string(rune(i))), Name: "Recipe"}
		}
		return recipes, nil
	}
	req2, _ := http.NewRequest("GET", "/recipes", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200 for large list, got %d", w2.Code)
	}

	// Special characters in tag
	repo.getByTagFunc = func(ctx context.Context, tag string) ([]model.Recipe, error) {
		return []model.Recipe{}, nil
	}
	req3, _ := http.NewRequest("GET", "/recipes/search?tag=tag%20with%20spaces", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Errorf("Expected status 200 for special chars, got %d", w3.Code)
	}
}

func stringPtr(s string) *string {
	return &s
}
