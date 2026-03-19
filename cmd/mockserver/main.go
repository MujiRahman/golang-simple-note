package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	r := gin.Default()

	r.POST("/register", handleRegister)
	r.POST("/login", handleLogin)

	// protected
	r.POST("/notes", authMiddleware(), handleCreateNote)
	r.GET("/notes", authMiddleware(), handleListNotes)
	r.GET("/notes/:id", authMiddleware(), handleGetNote)
	r.PUT("/notes/:id", authMiddleware(), handleUpdateNote)
	r.DELETE("/notes/:id", authMiddleware(), handleDeleteNote)

	r.Run(":8081")
}

func handleRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if _, ok := users[req.Username]; ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		return
	}
	u := &User{ID: nextUser, Username: req.Username}
	users[req.Username] = u
	userPw[req.Username] = req.Password
	nextUser++
	c.JSON(http.StatusCreated, u)
}

func handleLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	u, ok := users[req.Username]
	if !ok || userPw[req.Username] != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	// mock token: "mock-<id>"
	token := "mock-" + strconv.FormatUint(uint64(u.ID), 10)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// authMiddleware: check header Authorization: Bearer mock-<id>
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ah := c.GetHeader("Authorization")
		if len(ah) < 8 || ah[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}
		token := ah[7:]
		// token format mock-<id>
		if len(token) < 6 || token[:5] != "mock-" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		idStr := token[5:]
		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// inject into context
		c.Set("user_id", uint(id64))
		c.Next()
	}
}

func getUserID(c *gin.Context) uint {
	if v, exists := c.Get("user_id"); exists {
		if uid, ok := v.(uint); ok {
			return uid
		}
	}
	return 0
}

func handleCreateNote(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	uid := getUserID(c)
	n := &Note{ID: nextNote, UserID: uid, Title: req.Title, Content: req.Content}
	notes[nextNote] = n
	nextNote++
	c.JSON(http.StatusCreated, n)
}

func handleListNotes(c *gin.Context) {
	uid := getUserID(c)
	var out []Note
	for _, n := range notes {
		if n.UserID == uid {
			out = append(out, *n)
		}
	}
	c.JSON(http.StatusOK, out)
}

func handleGetNote(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	n, ok := notes[uint(id64)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	uid := getUserID(c)
	if n.UserID != uid {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found or access denied"})
		return
	}
	c.JSON(http.StatusOK, n)
}

func handleUpdateNote(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	n, ok := notes[uint(id64)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	uid := getUserID(c)
	if n.UserID != uid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not found or access denied"})
		return
	}
	n.Title = req.Title
	n.Content = req.Content
	c.JSON(http.StatusOK, n)
}

func handleDeleteNote(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	n, ok := notes[uint(id64)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	uid := getUserID(c)
	if n.UserID != uid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not found or access denied"})
		return
	}
	delete(notes, uint(id64))
	c.JSON(http.StatusNoContent, nil)
}
