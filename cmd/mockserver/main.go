package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type Note struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	users         = map[string]*User{}  // username -> user
	userPw        = map[string]string{} // username -> password (plain for mock)
	notes         = map[uint]*Note{}
	nextUser uint = 1
	nextNote uint = 1
)

func main() {
	r := httprouter.New()

	r.POST("/register", handleRegister)
	r.POST("/login", handleLogin)

	// protected
	r.POST("/notes", auth(handleCreateNote))
	r.GET("/notes", auth(handleListNotes))
	r.GET("/notes/:id", auth(handleGetNote))
	r.PUT("/notes/:id", auth(handleUpdateNote))
	r.DELETE("/notes/:id", auth(handleDeleteNote))

	http.ListenAndServe(":8081", r)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

func handleRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if _, ok := users[req.Username]; ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username already used"})
		return
	}
	u := &User{ID: nextUser, Username: req.Username}
	users[req.Username] = u
	userPw[req.Username] = req.Password
	nextUser++
	writeJSON(w, http.StatusCreated, u)
}

func handleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	u, ok := users[req.Username]
	if !ok || userPw[req.Username] != req.Password {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	// mock token: "mock-<id>"
	token := "mock-" + strconv.FormatUint(uint64(u.ID), 10)
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// auth middleware: check header Authorization: Bearer mock-<id>
func auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ah := r.Header.Get("Authorization")
		if len(ah) < 8 || ah[:7] != "Bearer " {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid authorization header"})
			return
		}
		token := ah[7:]
		// token format mock-<id>
		if len(token) < 6 || token[:5] != "mock-" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}
		idStr := token[5:]
		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}
		// inject into context
		ctx := context.WithValue(r.Context(), "user_id", uint(id64))
		h(w, r.WithContext(ctx), ps)
	}
}

func getUserID(r *http.Request) uint {
	if v := r.Context().Value("user_id"); v != nil {
		if uid, ok := v.(uint); ok {
			return uid
		}
	}
	return 0
}

func handleCreateNote(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req struct{ Title, Content string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	uid := getUserID(r)
	n := &Note{ID: nextNote, UserID: uid, Title: req.Title, Content: req.Content}
	notes[nextNote] = n
	nextNote++
	writeJSON(w, http.StatusCreated, n)
}

func handleListNotes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uid := getUserID(r)
	var out []Note
	for _, n := range notes {
		if n.UserID == uid {
			out = append(out, *n)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func handleGetNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	n, ok := notes[uint(id64)]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	uid := getUserID(r)
	if n.UserID != uid {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found or access denied"})
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func handleUpdateNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var req struct{ Title, Content string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	n, ok := notes[uint(id64)]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	uid := getUserID(r)
	if n.UserID != uid {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "not found or access denied"})
		return
	}
	n.Title = req.Title
	n.Content = req.Content
	writeJSON(w, http.StatusOK, n)
}

func handleDeleteNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	n, ok := notes[uint(id64)]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	uid := getUserID(r)
	if n.UserID != uid {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "not found or access denied"})
		return
	}
	delete(notes, uint(id64))
	writeJSON(w, http.StatusNoContent, nil)
}
