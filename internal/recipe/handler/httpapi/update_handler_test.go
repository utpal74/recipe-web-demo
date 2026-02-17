package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"github.com/gin-gonic/gin"
)

// Ensure UpdateRecipeHandler passes URI id into the domain update input
func TestUpdateRecipeHandler_PassesIDToService(t *testing.T) {
    gin.SetMode(gin.TestMode)

    var receivedID domain.RecipeID
    var receivedInput domain.UpdateRecipeInput

    // mock service
    mock := &mockService{}
    mock.updateFunc = func(ctx context.Context, input domain.UpdateRecipeInput) (domain.Recipe, error) {
        receivedID = input.ID
        receivedInput = input
        return domain.Recipe{ID: input.ID, Name: "ok"}, nil
    }

    handler := New(mock)
    r := gin.New()
    r.PUT("/recipes/:id", handler.UpdateRecipeHandler)

    body := updateRecipeRequest{Name: stringPtr("updated name")}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("PUT", "/recipes/abc123", bytes.NewBuffer(b))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }

    if receivedID != domain.RecipeID("abc123") {
        t.Fatalf("expected service update called with ID 'abc123', got '%s'", receivedID)
    }

    if receivedInput.Name == nil || *receivedInput.Name != "updated name" {
        t.Fatalf("expected name 'updated name' in update input, got %#v", receivedInput.Name)
    }
}
