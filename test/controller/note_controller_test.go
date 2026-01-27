package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/app"
	"github.com/MujiRahman/golang-simple-note/internal/model"
)

// fakeUserService for middleware token parsing
type fakeUserSvcForAuth struct{}

func (f *fakeUserSvcForAuth) Register(username, password string) (*model.User, error) {
	return nil, nil
}
func (f *fakeUserSvcForAuth) Login(username, password string) (string, error) { return "tok-1", nil }
func (f *fakeUserSvcForAuth) ParseToken(tokenStr string) (uint, error) {
	if tokenStr == "tok-1" {
		return 7, nil
	}
	return 0, nil
}

// fake NoteService that records created notes
type fakeNoteSvc struct {
	created *model.Note
}

func (f *fakeNoteSvc) Create(userID uint, title, content string) (*model.Note, error) {
	n := &model.Note{ID: 11, UserID: userID, Title: title, Content: content}
	f.created = n
	return n, nil
}
func (f *fakeNoteSvc) GetByID(userID, id uint) (*model.Note, error) { return f.created, nil }
func (f *fakeNoteSvc) ListByUser(userID uint) ([]model.Note, error) {
	return []model.Note{*f.created}, nil
}
func (f *fakeNoteSvc) Update(userID, id uint, title, content string) (*model.Note, error) {
	return f.created, nil
}
func (f *fakeNoteSvc) Delete(userID, id uint) error { return nil }

func TestCreateNote_Unauthorized(t *testing.T) {
	us := &fakeUserSvcForAuth{}
	ns := &fakeNoteSvc{}
	router := app.NewRouter(us, ns, &config.Config{})

	body := map[string]string{"title": "t1", "content": "c1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateNote_Success(t *testing.T) {
	us := &fakeUserSvcForAuth{}
	ns := &fakeNoteSvc{}
	router := app.NewRouter(us, ns, &config.Config{})

	body := map[string]string{"title": "t1", "content": "c1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer tok-1")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
