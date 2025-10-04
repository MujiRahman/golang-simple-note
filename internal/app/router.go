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
	r.POST("/register", wrapHandler(userCtrl.Register))
	r.POST("/login", wrapHandler(userCtrl.Login))

	// protected group: since httprouter doesn't have groups, attach middleware per route
	// We'll use Auth middleware wrapper that needs user service to parse token.
	authMw := middleware.AuthMiddleware(userSvc)

	r.POST("/notes", authMw(wrapHandler(noteCtrl.Create)))
	r.GET("/notes", authMw(wrapHandler(noteCtrl.List)))
	r.GET("/notes/:id", authMw(wrapHandler(noteCtrl.Get)))
	r.PUT("/notes/:id", authMw(wrapHandler(noteCtrl.Update)))
	r.DELETE("/notes/:id", authMw(wrapHandler(noteCtrl.Delete)))

	// fallback
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("noteapp up"))
	})

	return r
}

// Convert httprouter.Handle to http.HandlerFunc using closure so we can use middleware that accepts http.HandlerFunc-like
func wrapHandler(h httprouter.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// but httprouter passes Params; we don't have them here.
		// This wrapper is only used when middleware expects http.HandlerFunc.
		// For routes registered directly with r.GET/POST/etc we pass a final adapter, but above we wrap again.
		// To avoid complexity, we will set Params via context in middleware if needed.
		// Simpler approach: call original with empty params and rely on controller to read from URL via httprouter.RequestContext
		// But httprouter stores params in request context under key "params". To keep simple, we will not use this wrapper for param routes.
		// However our router registration above used wrapHandler consistently — httprouter will set params when routing.
		// The wrapped handler needs to recover params from the request context.
		// Use httprouter.ParamsFromContext
		ps := httprouter.ParamsFromContext(r.Context())
		h(w, r, ps)
	}
}
