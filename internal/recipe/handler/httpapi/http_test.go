package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"github.com/gin-gonic/gin"
)

// Mock service for testing
type mockService struct {
	createFunc     func(context.Context, domain.CreateRecipeInput) (domain.Recipe, error)
	getByIDFunc    func(context.Context, domain.RecipeID) (domain.Recipe, error)
	getAllFunc     func(context.Context) ([]domain.Recipe, error)
	updateFunc     func(context.Context, domain.UpdateRecipeInput) (domain.Recipe, error)
	deleteFunc     func(context.Context, domain.RecipeID) error
	getByTagFunc   func(context.Context, string) ([]domain.Recipe, error)
	createManyFunc func(context.Context, []domain.Recipe) error
}

func (m *mockService) Create(ctx context.Context, input domain.CreateRecipeInput) (domain.Recipe, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return domain.Recipe{
		ID:           domain.RecipeID("test-id"),
		Name:         input.Name,
		Tags:         input.Tags,
		Ingredients:  input.Ingredients,
		Instructions: input.Instructions,
	}, nil
}

func (m *mockService) GetByID(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return domain.Recipe{ID: id, Name: "Test"}, nil
}

func (m *mockService) GetAll(ctx context.Context) ([]domain.Recipe, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return []domain.Recipe{}, nil
}

func (m *mockService) Update(ctx context.Context, input domain.UpdateRecipeInput) (domain.Recipe, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, input)
	}
	return domain.Recipe{ID: input.ID}, nil
}

func (m *mockService) Delete(ctx context.Context, id domain.RecipeID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockService) GetByTag(ctx context.Context, tag string) ([]domain.Recipe, error) {
	if m.getByTagFunc != nil {
		return m.getByTagFunc(ctx, tag)
	}
	return []domain.Recipe{}, nil
}

func (m *mockService) CreateMany(ctx context.Context, recipes []domain.Recipe) error {
	if m.createManyFunc != nil {
		return m.createManyFunc(ctx, recipes)
	}
	return nil
}

// Helper to setup router with handlers
func setupTestRouter(service *mockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := New(service)

	router.GET("/recipes", handler.ListRecipeHandler)
	router.GET("/recipes/search", handler.ListRecipesByTagHandler)
	router.GET("/recipes/:id", handler.GetRecipeByIDHandler)
	router.POST("/recipes", handler.CreateRecipeHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)

	return router
}

func TestCreateRecipeHandler(t *testing.T) {
	service := &mockService{}
	router := setupTestRouter(service)

	recipe := createRecipeRequest{Name: "Test Recipe", Tags: []string{"tag1"}, Ingredients: []string{"ing"}, Instructions: []string{"step"}}
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
	service.createFunc = func(ctx context.Context, input domain.CreateRecipeInput) (domain.Recipe, error) {
		return domain.Recipe{}, errors.New("persistence error")
	}
	req3, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer(body))
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w3.Code)
	}
}

func TestListRecipeHandler(t *testing.T) {
	service := &mockService{
		getAllFunc: func(ctx context.Context) ([]domain.Recipe, error) {
			return []domain.Recipe{{ID: "1", Name: "R1"}, {ID: "2", Name: "R2"}}, nil
		},
	}
	router := setupTestRouter(service)

	req, _ := http.NewRequest("GET", "/recipes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Error case
	service.getAllFunc = func(ctx context.Context) ([]domain.Recipe, error) {
		return nil, errors.New("error")
	}
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)

	if w2.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w2.Code)
	}
}

func TestUpdateRecipeHandler(t *testing.T) {
	service := &mockService{}
	router := setupTestRouter(service)

	cmd := updateRecipeRequest{Name: stringPtr("Updated")}
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
	service.updateFunc = func(ctx context.Context, input domain.UpdateRecipeInput) (domain.Recipe, error) {
		return domain.Recipe{}, domain.ErrNotFound
	}
	req3, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer(body))
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w3.Code)
	}

	// Persistence error
	service.updateFunc = func(ctx context.Context, input domain.UpdateRecipeInput) (domain.Recipe, error) {
		return domain.Recipe{}, domain.ErrPersistence
	}
	req4, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer(body))
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w4.Code)
	}
}

func TestListRecipesByTagHandler(t *testing.T) {
	service := &mockService{
		getByTagFunc: func(ctx context.Context, tag string) ([]domain.Recipe, error) {
			if tag == "test" {
				return []domain.Recipe{{ID: "1", Tags: []string{"test"}}}, nil
			}
			return []domain.Recipe{}, nil
		},
	}
	router := setupTestRouter(service)

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
	service.getByTagFunc = func(ctx context.Context, tag string) ([]domain.Recipe, error) {
		return nil, domain.ErrInvalidInput
	}
	req3, _ := http.NewRequest("GET", "/recipes/search?tag=empty", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w3.Code)
	}
}

func TestGetRecipeByIDHandler(t *testing.T) {
	service := &mockService{
		getByIDFunc: func(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
			if id == "1" {
				return domain.Recipe{ID: "1", Name: "Test"}, nil
			}
			return domain.Recipe{}, domain.ErrNotFound
		},
	}
	router := setupTestRouter(service)

	req, _ := http.NewRequest("GET", "/recipes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Not found
	service.getByIDFunc = func(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
		return domain.Recipe{}, domain.ErrNotFound
	}
	req2, _ := http.NewRequest("GET", "/recipes/nonexistent", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w2.Code)
	}

	// Persistence error
	service.getByIDFunc = func(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
		return domain.Recipe{}, errors.New("database error")
	}
	req3, _ := http.NewRequest("GET", "/recipes/1", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w3.Code)
	}
}

func TestDeleteRecipeHandler(t *testing.T) {
	service := &mockService{
		deleteFunc: func(ctx context.Context, id domain.RecipeID) error {
			if id == "1" {
				return nil
			}
			return domain.ErrNotFound
		},
	}
	router := setupTestRouter(service)

	req, _ := http.NewRequest("DELETE", "/recipes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Not found
	service.deleteFunc = func(ctx context.Context, id domain.RecipeID) error {
		return domain.ErrNotFound
	}
	req2, _ := http.NewRequest("DELETE", "/recipes/nonexistent", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w2.Code)
	}

	// Persistence error
	service.deleteFunc = func(ctx context.Context, id domain.RecipeID) error {
		return errors.New("database error")
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
	service := &mockService{}
	router := setupTestRouter(service)

	// Minimal valid recipe
	recipe := createRecipeRequest{
		Name:         "Minimal",
		Tags:         []string{"tag1"},
		Ingredients:  []string{"ing1"},
		Instructions: []string{"inst1"},
	}
	body, _ := json.Marshal(recipe)
	req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Should create with minimal valid data
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 for minimal valid data, got %d", w.Code)
	}

	// Large lists
	service.getAllFunc = func(ctx context.Context) ([]domain.Recipe, error) {
		recipes := make([]domain.Recipe, 1000)
		for i := 0; i < 1000; i++ {
			recipes[i] = domain.Recipe{ID: domain.RecipeID("id" + string(rune(i))), Name: "Recipe"}
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
	service.getByTagFunc = func(ctx context.Context, tag string) ([]domain.Recipe, error) {
		return []domain.Recipe{}, nil
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
