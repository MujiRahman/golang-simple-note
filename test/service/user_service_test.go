package service_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/model"
	"github.com/MujiRahman/golang-simple-note/internal/service"
)

// fake user repo implementing repository.UserRepository (minimal for tests)
type fakeUserRepo struct {
	users map[string]*model.User
}

func (f *fakeUserRepo) Create(u *model.User) error {
	if f.users == nil {
		f.users = map[string]*model.User{}
	}
	f.users[u.Username] = u
	return nil
}
func (f *fakeUserRepo) FindByUsername(username string) (*model.User, error) {
	if u, ok := f.users[username]; ok {
		return u, nil
	}
	return nil, nil
}
func (f *fakeUserRepo) FindByID(id uint) (*model.User, error) { return nil, nil }

func TestUserService_RegisterLogin_ParseToken(t *testing.T) {
	repo := &fakeUserRepo{users: map[string]*model.User{}}
	cfg := &config.Config{JWTSecret: "testsecret", TokenTTL: 3600}
	svc := service.NewUserService(repo, cfg)

	// register
	u, err := svc.Register("alice", "password123")
	if err != nil {
		t.Fatalf("register error: %v", err)
	}
	if u.Username != "alice" {
		t.Fatalf("unexpected username: %v", u.Username)
	}

	// prepare stored user with bcrypt hash for login
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo.users["bob"] = &model.User{ID: 42, Username: "bob", Password: string(hash)}

	tok, err := svc.Login("bob", "password123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if tok == "" {
		t.Fatalf("empty token")
	}

	uid, err := svc.ParseToken(tok)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}
	if uid != 42 {
		t.Fatalf("expected uid 42, got %d", uid)
	}
}
