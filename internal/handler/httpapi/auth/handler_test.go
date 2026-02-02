package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-demo/recipes-web/internal/handler/httpapi/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestSignInHandler_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    ah := New(Config{Secret: "test-secret", Issuer: "test-issuer"})
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

    var out JWTOutput
    if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
        t.Fatalf("failed to unmarshal response: %v", err)
    }
    if out.Token == "" {
        t.Fatalf("expected token in response")
    }

    claims := Claims{}
    token, err := jwt.ParseWithClaims(out.Token, &claims, func(t *jwt.Token) (any, error) {
        if t.Method != jwt.SigningMethodHS256 {
            return nil, nil
        }
        return []byte("test-secret"), nil
    })
    if err != nil || !token.Valid {
        t.Fatalf("token invalid: %v", err)
    }
    if claims.UserName != "admin" {
        t.Errorf("unexpected userName claim: %s", claims.UserName)
    }
    if claims.Role != "admin" {
        t.Errorf("unexpected role claim: %s", claims.Role)
    }
}

func TestSignInHandler_BadCredentialsAndBadJSON(t *testing.T) {
    gin.SetMode(gin.TestMode)
    ah := New(Config{Secret: "test-secret", Issuer: "test-issuer"})
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

    // configure secret for middleware
    secret := "middleware-secret"
    os.Setenv("JWT_SECRET", secret)
    defer os.Unsetenv("JWT_SECRET")

    ah := New(Config{Secret: secret, Issuer: "test-issuer"})
    expiry := time.Now().Add(10 * time.Minute)
    token, err := ah.createToken("admin", expiry)
    if err != nil {
        t.Fatalf("failed to create token: %v", err)
    }

    router := gin.New()
    router.GET("/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
        c.JSON(200, gin.H{"userName": c.GetString("userName"), "role": c.GetString("role")})
    })

    // Valid token
    req, _ := http.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200 with valid token, got %d, body: %s", w.Code, w.Body.String())
    }

    // Missing/invalid Authorization header
    req2, _ := http.NewRequest("GET", "/protected", nil)
    w2 := httptest.NewRecorder()
    router.ServeHTTP(w2, req2)
    if w2.Code != http.StatusUnauthorized {
        t.Errorf("expected 401 when authorization header missing, got %d", w2.Code)
    }
}
