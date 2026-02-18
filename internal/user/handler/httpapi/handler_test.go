package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-demo/recipes-web/internal/user/domain"
	"github.com/gin-gonic/gin"
)

func TestUpdateUserHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{}
	handler := New(mockService)
	router := gin.New()
	router.PUT("/users/:id", handler.UpdateUserHandler)

	body := `{"userName":"updated_name"}`
	req, _ := http.NewRequest("PUT", "/users/user123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Check for userName field (lowercase due to JSON marshaling)
	if val, ok := response["UserName"]; !ok || val != "updated_name" {
		t.Errorf("expected UserName 'updated_name', got %v", response["UserName"])
	}
}

func TestUpdateUserHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{}
	handler := New(mockService)
	router := gin.New()
	router.PUT("/users/:id", handler.UpdateUserHandler)

	req, _ := http.NewRequest("PUT", "/users/user123", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateUserHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockUserService{}
	h := New(mockService)
	router := gin.New()
	router.POST("/update", h.UpdateUserHandler)

	req, _ := http.NewRequest("POST", "/update", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid json, got %d", w.Code)
	}
}

func TestUpdateUserHandler_Persistence_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{
		updateErr: domain.ErrPersistence,
	}
	handler := New(mockService)
	router := gin.New()
	router.PUT("/users/:id", handler.UpdateUserHandler)

	body := `{"userName":"updated_name"}`
	req, _ := http.NewRequest("PUT", "/users/user123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotModified {
		t.Errorf("expected status 304, got %d", w.Code)
	}
}

func TestUpdateUserHandler_Internal_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{
		updateErr: errors.New("some internal error"),
	}
	handler := New(mockService)
	router := gin.New()
	router.PUT("/users/:id", handler.UpdateUserHandler)

	body := `{"userName":"updated_name"}`
	req, _ := http.NewRequest("PUT", "/users/user123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestDeleteUserHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{}
	handler := New(mockService)
	router := gin.New()
	router.DELETE("/users/:id", handler.DeleteUserHandler)

	req, _ := http.NewRequest("DELETE", "/users/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestDeleteUserHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{
		deleteErr: domain.ErrNotFound,
	}
	handler := New(mockService)
	router := gin.New()
	router.DELETE("/users/:id", handler.DeleteUserHandler)

	req, _ := http.NewRequest("DELETE", "/users/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteUserHandler_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{
		deleteErr: errors.New("internal error"),
	}
	handler := New(mockService)
	router := gin.New()
	router.DELETE("/users/:id", handler.DeleteUserHandler)

	req, _ := http.NewRequest("DELETE", "/users/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestFindUserByIDHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{}
	handler := New(mockService)
	router := gin.New()
	router.GET("/users/:id", handler.FindUserByIDHandler)

	req, _ := http.NewRequest("GET", "/users/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if val, ok := response["UserName"]; !ok || val != "test_user" {
		t.Errorf("expected UserName 'test_user', got %v", response["UserName"])
	}
}

func TestFindUserByIDHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{
		findByIDErr: domain.ErrNotFound,
	}
	handler := New(mockService)
	router := gin.New()
	router.GET("/users/:id", handler.FindUserByIDHandler)

	req, _ := http.NewRequest("GET", "/users/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestFindUserByIDHandler_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{
		findByIDErr: errors.New("internal error"),
	}
	handler := New(mockService)
	router := gin.New()
	router.GET("/users/:id", handler.FindUserByIDHandler)

	req, _ := http.NewRequest("GET", "/users/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestFindUserByNameHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{}
	handler := New(mockService)
	router := gin.New()
	router.GET("/users/name/:name", handler.FindUserByNameHandler)

	req, _ := http.NewRequest("GET", "/users/name/john_doe", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if val, ok := response["UserName"]; !ok || val != "john_doe" {
		t.Errorf("expected UserName 'john_doe', got %v", response["UserName"])
	}
}

func TestFindUserByNameHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{
		findByNameErr: domain.ErrNotFound,
	}
	handler := New(mockService)
	router := gin.New()
	router.GET("/users/name/:name", handler.FindUserByNameHandler)

	req, _ := http.NewRequest("GET", "/users/name/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestFindUserByNameHandler_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockUserService{
		findByNameErr: errors.New("internal error"),
	}
	handler := New(mockService)
	router := gin.New()
	router.GET("/users/name/:name", handler.FindUserByNameHandler)

	req, _ := http.NewRequest("GET", "/users/name/john_doe", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

// Mock UserService for testing
type mockUserService struct {
	createErr     error
	updateErr     error
	deleteErr     error
	findByIDErr   error
	findByNameErr error
}

func (m *mockUserService) Create(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	if m.createErr != nil {
		return domain.User{}, m.createErr
	}
	return domain.User{
		ID:           domain.UserID("user123"),
		UserName:     input.UserName,
		PasswordHash: "hashed_password",
		Role:         domain.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *mockUserService) Update(ctx context.Context, input domain.UpdateUserInput) (domain.User, error) {
	if m.updateErr != nil {
		return domain.User{}, m.updateErr
	}
	userName := "updated_name"
	if input.UserName != nil {
		userName = *input.UserName
	}
	return domain.User{
		ID:           input.ID,
		UserName:     userName,
		PasswordHash: "hashed_password",
		Role:         domain.RoleUser,
		CreatedAt:    time.Now().Add(-time.Hour),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *mockUserService) Delete(ctx context.Context, id domain.UserID) error {
	return m.deleteErr
}

func (m *mockUserService) FindByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	if m.findByIDErr != nil {
		return domain.User{}, m.findByIDErr
	}
	return domain.User{
		ID:           id,
		UserName:     "test_user",
		PasswordHash: "hashed_password",
		Role:         domain.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *mockUserService) FindByUserName(ctx context.Context, name string) (domain.User, error) {
	if m.findByNameErr != nil {
		return domain.User{}, m.findByNameErr
	}
	return domain.User{
		ID:           domain.UserID("user123"),
		UserName:     name,
		PasswordHash: "hashed_password",
		Role:         domain.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
