package app

import (
	"net/http"
	// "strings"

	"github.com/julienschmidt/httprouter"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/controller"

	// "github.com/MujiRahman/golang-simple-note/internal/repository"
	"github.com/MujiRahman/golang-simple-note/internal/service"
	"github.com/MujiRahman/golang-simple-note/pkg/middleware"
)

// NewRouter builds router with DI
func NewRouter(userSvc service.UserService, noteSvc service.NoteService, cfg *config.Config) http.Handler {
	r := httprouter.New()

	// controllers
	userCtrl := controller.NewUserController(userSvc)
	noteCtrl := controller.NewNoteController(noteSvc)

	// public
	r.POST("/register", userCtrl.Register)
	r.POST("/login", userCtrl.Login)

	// protected group: since httprouter doesn't have groups, attach middleware per route
	// We'll use Auth middleware wrapper that needs user service to parse token.
	authMw := middleware.AuthMiddleware(userSvc)

	r.POST("/notes", authMw(noteCtrl.Create))
	r.GET("/notes", authMw(noteCtrl.List))
	r.GET("/notes/:id", authMw(noteCtrl.Get))
	r.PUT("/notes/:id", authMw(noteCtrl.Update))
	r.DELETE("/notes/:id", authMw(noteCtrl.Delete))

	// fallback
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("noteapp up"))
	})

	return r
}

// No adapter needed: controllers and middleware use httprouter.Handle directly.
