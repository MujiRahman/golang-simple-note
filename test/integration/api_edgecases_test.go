package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Register with invalid JSON should return 400
func TestE2E_RegisterInvalidPayload(t *testing.T) {
	router := setupRouterForTest(t)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewReader([]byte("not-json")))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 bad request, got %d", resp.StatusCode)
	}
}

// Login with wrong password should return 401
func TestE2E_LoginInvalidCredentials(t *testing.T) {
	router := setupRouterForTest(t)
	server := httptest.NewServer(router)
	defer server.Close()

	// register first
	regBody := map[string]string{"username": "e2euser2", "password": "pass"}
	rb, _ := json.Marshal(regBody)
	resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewReader(rb))
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d", resp.StatusCode)
	}

	// try login with wrong password
	bad := map[string]string{"username": "e2euser2", "password": "wrong"}
	bb, _ := json.Marshal(bad)
	resp, err = http.Post(server.URL+"/login", "application/json", bytes.NewReader(bb))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 unauthorized for wrong credentials, got %d", resp.StatusCode)
	}
}

// Create note without Authorization header should return 401
func TestE2E_CreateNote_Unauthorized(t *testing.T) {
	router := setupRouterForTest(t)
	server := httptest.NewServer(router)
	defer server.Close()

	noteBody := map[string]string{"title": "t1", "content": "c1"}
	nb, _ := json.Marshal(noteBody)
	resp, err := http.Post(server.URL+"/notes", "application/json", bytes.NewReader(nb))
	if err != nil {
		t.Fatalf("create note request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 when creating note without token, got %d", resp.StatusCode)
	}
}

// Create note with invalid JSON payload should return 400 even with valid token
func TestE2E_CreateNote_InvalidPayload(t *testing.T) {
	router := setupRouterForTest(t)
	server := httptest.NewServer(router)
	defer server.Close()

	// register & login
	regBody := map[string]string{"username": "e2euser3", "password": "pass"}
	rb, _ := json.Marshal(regBody)
	resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewReader(rb))
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d", resp.StatusCode)
	}

	resp, err = http.Post(server.URL+"/login", "application/json", bytes.NewReader(rb))
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 ok on login, got %d", resp.StatusCode)
	}
	var loginResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login resp: %v", err)
	}
	token := loginResp["token"]

	// send malformed json
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/notes", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid note payload, got %d", resp.StatusCode)
	}
}
