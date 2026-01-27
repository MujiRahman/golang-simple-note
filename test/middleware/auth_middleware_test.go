package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"

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

func protectedHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := r.Context().Value(contextkey.UserIDKey)
	if v == nil {
		http.Error(w, "no user", http.StatusInternalServerError)
		return
	}
	uid := v.(uint)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%d", uid)))
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	mw := middleware.AuthMiddleware(&fakeUserService{uid: 0, err: nil})
	h := mw(protectedHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h(rr, req, httprouter.Params{})
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	mw := middleware.AuthMiddleware(&fakeUserService{uid: 0, err: nil})
	h := mw(protectedHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bad token")
	rr := httptest.NewRecorder()
	h(rr, req, httprouter.Params{})
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mw := middleware.AuthMiddleware(&fakeUserService{uid: 0, err: fmt.Errorf("invalid")})
	h := mw(protectedHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	rr := httptest.NewRecorder()
	h(rr, req, httprouter.Params{})
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mw := middleware.AuthMiddleware(&fakeUserService{uid: 77, err: nil})
	h := mw(protectedHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer goodtoken")
	rr := httptest.NewRecorder()
	h(rr, req, httprouter.Params{})
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if rr.Body.String() != "77" {
		t.Fatalf("expected body 77, got %s", rr.Body.String())
	}
}
