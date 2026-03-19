package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/MujiRahman/golang-simple-note/internal/model"
	"github.com/MujiRahman/golang-simple-note/pkg/contextkey"
	"github.com/MujiRahman/golang-simple-note/pkg/middleware"
)

// fakeUserService implements the minimal UserService methods for middleware tests
type fakeUserService struct {
	uid uint
	err error
}

func (f *fakeUserService) Register(username, password string) (*model.User, error) {
	return &model.User{ID: f.uid, Username: username}, nil
}
func (f *fakeUserService) Login(username, password string) (string, error) { return "", nil }
func (f *fakeUserService) ParseToken(tokenStr string) (uint, error)        { return f.uid, f.err }

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	mw := middleware.AuthMiddleware(&fakeUserService{uid: 0, err: nil})

	// Create a test gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	// Call the middleware
	mw(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	mw := middleware.AuthMiddleware(&fakeUserService{uid: 0, err: nil})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bad token")
	c.Request = req

	mw(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mw := middleware.AuthMiddleware(&fakeUserService{uid: 0, err: fmt.Errorf("invalid")})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	c.Request = req

	mw(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mw := middleware.AuthMiddleware(&fakeUserService{uid: 77, err: nil})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer goodtoken")
	c.Request = req

	// Call middleware
	mw(c)

	// Check that the token was injected
	uid := c.GetUint(string(contextkey.UserIDKey))
	if uid != 77 {
		t.Fatalf("expected uid 77, got %d", uid)
	}

	if w.Code != http.StatusOK && w.Code != 0 {
		t.Fatalf("expected 200 or no status (middleware passed through), got %d", w.Code)
	}
}
