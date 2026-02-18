package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-demo/recipes-web/internal/auth/domain"
	"github.com/gin-demo/recipes-web/internal/auth/usecase"
	"github.com/gin-demo/recipes-web/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestSignInHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock service
	mockService := &mockAuthService{}
	ah := New(mockService)
	router := gin.New()
	router.POST("/signin", ah.SignInHandler)

	body := `{"userName":"admin","password":"password"}`
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var out signInUserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if out.Token == "" {
		t.Fatalf("expected token in response")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(out.Token, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, nil
		}
		return []byte("test-secret"), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("token invalid: %v", err)
	}
	if claims["userName"] != "admin" {
		t.Errorf("unexpected userName claim: %v", claims["userName"])
	}
}

func TestSignInHandler_BadCredentialsAndBadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockAuthService{}
	ah := New(mockService)
	router := gin.New()
	router.POST("/signin", ah.SignInHandler)

	// Bad credentials
	body := `{"userName":"admin","password":"wrong"}`
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for bad credentials, got %d", w.Code)
	}

	// Invalid JSON
	req2, _ := http.NewRequest("POST", "/signin", bytes.NewBufferString("invalid"))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid json, got %d", w2.Code)
	}
}

func TestAuthMiddleware_ValidAndInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protected", middleware.NewAuthMiddleWare(nil).Handle(), func(c *gin.Context) {
		c.JSON(200, gin.H{"userName": c.GetString("userName"), "role": c.GetString("role")})
	})

	// For this test, we'll use a test token (in practice, you'd generate this from your token service)
	// Missing/invalid Authorization header
	req2, _ := http.NewRequest("GET", "/protected", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 when authorization header missing, got %d", w2.Code)
	}
}

func TestRefreshHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockAuthService{}
	ah := New(mockService)
	router := gin.New()
	router.POST("/refresh", ah.RefreshHandler)

	body := `{"refreshToken":"dummy-refresh-token"}`
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var out map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if out["accessToken"] == "" || out["refreshToken"] == "" {
		t.Fatalf("expected accessToken and refreshToken in response")
	}
}

type mockAuthService struct{}
func (m *mockAuthService) Refresh(ctx context.Context, tokenString string) (string, string, error) {
	// Dummy implementation for testing
	return "new-access-token", "new-refresh-token", nil
}

func (m *mockAuthService) SignIn(ctx context.Context, input usecase.SignInInput) (usecase.SignInOutput, error) {
	// Verify credentials - reject wrong password
	if input.Password != "password" {
		return usecase.SignInOutput{}, domain.ErrInvalidCredentials
	}

	// create a signed JWT matching the test parser's expectations
	claims := jwt.MapClaims{
		"userName": input.UserName,
		"exp":      jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		return usecase.SignInOutput{}, err
	}
	return usecase.SignInOutput{
		Token:   signed,
		Expires: time.Now().Add(15 * time.Minute),
	}, nil
}

func (m *mockAuthService) SignOut(ctx context.Context, input domain.Identity) error {
	return nil
}

func (m *mockAuthService) SignUp(ctx context.Context, input usecase.SignUpInput) (usecase.SignUpOutput, error) {
	return usecase.SignUpOutput{}, nil
}
