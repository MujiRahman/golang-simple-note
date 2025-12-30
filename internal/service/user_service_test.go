package service
package service

import (
    "testing"

    "golang.org/x/crypto/bcrypt"

    "github.com/MujiRahman/golang-simple-note/config"
    "github.com/MujiRahman/golang-simple-note/internal/model"
)

type mockUserRepo struct {
    byName map[string]*model.User
    byID   map[uint]*model.User
    nextID uint
}

func newMockUserRepo() *mockUserRepo {
    return &mockUserRepo{byName: make(map[string]*model.User), byID: make(map[uint]*model.User), nextID: 1}
}

func (m *mockUserRepo) Create(user *model.User) error {
    if user.ID == 0 {
        user.ID = m.nextID
        m.nextID++
    }
    // store
    m.byName[user.Username] = user
    m.byID[user.ID] = user
    return nil
}

func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
    u, ok := m.byName[username]
    if !ok {
        return nil, nil
    }
    return u, nil
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
    u, ok := m.byID[id]
    if !ok {
        return nil, nil
    }
    return u, nil
}

func TestUserService_RegisterAndLogin_ParseToken(t *testing.T) {
    repo := newMockUserRepo()
    cfg := &config.Config{JWTSecret: "testsecret", TokenTTL: 3600}
    svc := NewUserService(repo, cfg)

    // Register
    u, err := svc.Register("alice", "password123")
    if err != nil {
        t.Fatalf("Register failed: %v", err)
    }
    if u.Username != "alice" {
        t.Fatalf("unexpected username: %v", u.Username)
    }
    // password should be hashed
    if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("password123")); err != nil {
        t.Fatalf("stored password is not a valid bcrypt hash: %v", err)
    }

    // Login
    token, err := svc.Login("alice", "password123")
    if err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    if token == "" {
        t.Fatalf("empty token returned")
    }

    // Parse token
    uid, err := svc.ParseToken(token)
    if err != nil {
        t.Fatalf("ParseToken failed: %v", err)
    }
    if uid != u.ID {
        t.Fatalf("unexpected uid from token: got %d want %d", uid, u.ID)
    }
}

func TestUserService_Register_Existing(t *testing.T) {
    repo := newMockUserRepo()
    cfg := &config.Config{JWTSecret: "s", TokenTTL: 3600}
    svc := NewUserService(repo, cfg)

    // Create existing user in repo
    existing := &model.User{Username: "bob", Password: "x"}
    repo.Create(existing)

    _, err := svc.Register("bob", "pw")
    if err == nil {
        t.Fatalf("expected error when registering existing username")
    }
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
    repo := newMockUserRepo()
    cfg := &config.Config{JWTSecret: "s", TokenTTL: 3600}
    svc := NewUserService(repo, cfg)

    // prepare user with hashed password
    hashed, _ := bcrypt.GenerateFromPassword([]byte("rightpw"), bcrypt.DefaultCost)
    user := &model.User{Username: "carol", Password: string(hashed)}
    repo.Create(user)

    _, err := svc.Login("carol", "wrongpw")
    if err == nil {
        t.Fatalf("expected error for wrong password")
    }
}
