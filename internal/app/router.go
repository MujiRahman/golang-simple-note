package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/controller"
	"github.com/MujiRahman/golang-simple-note/internal/service"
	"github.com/MujiRahman/golang-simple-note/pkg/middleware"
)

// NewRouter builds router with DI
func NewRouter(userSvc service.UserService, noteSvc service.NoteService, cfg *config.Config) http.Handler {
	r := gin.Default()

	// controllers
	userCtrl := controller.NewUserController(userSvc)
	noteCtrl := controller.NewNoteController(noteSvc)

	// public
	r.POST("/register", userCtrl.Register)
	r.POST("/login", userCtrl.Login)

	// protected group: using gin middleware
	authMw := middleware.AuthMiddleware(userSvc)

	r.POST("/notes", authMw, noteCtrl.Create)
	r.GET("/notes", authMw, noteCtrl.List)
	r.GET("/notes/:id", authMw, noteCtrl.Get)
	r.PUT("/notes/:id", authMw, noteCtrl.Update)
	r.DELETE("/notes/:id", authMw, noteCtrl.Delete)

	// fallback
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "noteapp up")
	})

	return r
}
