package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/app"
	"github.com/MujiRahman/golang-simple-note/internal/model"
	"github.com/MujiRahman/golang-simple-note/internal/repository"
	"github.com/MujiRahman/golang-simple-note/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouterForTest(t *testing.T) http.Handler {
	t.Helper()
	// in-memory sqlite
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm sqlite: %v", err)
	}
	// migrate
	if err := gdb.AutoMigrate(&model.User{}, &model.Note{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	userRepo := repository.NewUserRepository(gdb)
	noteRepo := repository.NewNoteRepository(gdb)

	cfg := &config.Config{JWTSecret: "integration-secret", TokenTTL: 3600}
	userSvc := service.NewUserService(userRepo, cfg)
	noteSvc := service.NewNoteService(noteRepo)

	return app.NewRouter(userSvc, noteSvc, cfg)
}

func TestEndToEnd_RegisterLoginCreateList(t *testing.T) {
	router := setupRouterForTest(t)
	server := httptest.NewServer(router)
	defer server.Close()

	// register
	regBody := map[string]string{"username": "e2euser", "password": "pass"}
	b, _ := json.Marshal(regBody)
	resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d", resp.StatusCode)
	}

	// login
	b, _ = json.Marshal(regBody)
	resp, err = http.Post(server.URL+"/login", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 ok on login, got %d", resp.StatusCode)
	}
	var loginResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login resp: %v", err)
	}
	token := loginResp["token"]
	if token == "" {
		t.Fatalf("empty token from login")
	}

	// create note
	noteBody := map[string]string{"title": "hello", "content": "world"}
	nb, _ := json.Marshal(noteBody)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/notes", bytes.NewReader(nb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create note request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 created for note, got %d", resp.StatusCode)
	}
	var created model.Note
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created note: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected created note id, got 0")
	}

	// list notes
	req2, _ := http.NewRequest(http.MethodGet, server.URL+"/notes", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("list notes request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on list, got %d", resp.StatusCode)
	}
	var notes []model.Note
	if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
		t.Fatalf("decode notes: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
}
