package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/app"
	"github.com/MujiRahman/golang-simple-note/internal/model"
)

// fakeUserService implements service.UserService for controller tests
type fakeUserService struct {
	registered map[string]*model.User
}

func newFakeUserService() *fakeUserService {
	return &fakeUserService{registered: map[string]*model.User{}}
}

func (f *fakeUserService) Register(username, password string) (*model.User, error) {
	u := &model.User{ID: 100, Username: username}
	f.registered[username] = u
	return u, nil
}
func (f *fakeUserService) Login(username, password string) (string, error) {
	if _, ok := f.registered[username]; !ok {
		return "", errors.New("invalid credentials")
	}
	return "tok-123", nil
}
func (f *fakeUserService) ParseToken(tokenStr string) (uint, error) {
	// token "tok-123" maps to user id 100
	if tokenStr == "tok-123" {
		return 100, nil
	}
	return 0, errors.New("invalid token")
}

func TestRegisterHandler(t *testing.T) {
	us := newFakeUserService()
	ns := &fakeNoteService{} // not used here
	router := app.NewRouter(us, ns, &config.Config{})

	body := map[string]string{"username": "alice", "password": "pw"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestLoginHandler(t *testing.T) {
	us := newFakeUserService()
	us.Register("bob", "pw")
	ns := &fakeNoteService{}
	router := app.NewRouter(us, ns, &config.Config{})

	body := map[string]string{"username": "bob", "password": "pw"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

// minimal fakeNoteService used to satisfy NewRouter; real tests for notes in other file
type fakeNoteService struct{}

func (f *fakeNoteService) Create(userID uint, title, content string) (*model.Note, error) {
	return nil, nil
}
func (f *fakeNoteService) GetByID(userID, id uint) (*model.Note, error) { return nil, nil }
func (f *fakeNoteService) ListByUser(userID uint) ([]model.Note, error) { return nil, nil }
func (f *fakeNoteService) Update(userID, id uint, title, content string) (*model.Note, error) {
	return nil, nil
}
func (f *fakeNoteService) Delete(userID, id uint) error { return nil }
