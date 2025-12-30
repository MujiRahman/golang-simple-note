package main

import (
	"net/http"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/app"
	"github.com/MujiRahman/golang-simple-note/internal/repository"
	"github.com/MujiRahman/golang-simple-note/internal/service"
	"github.com/MujiRahman/golang-simple-note/pkg/logger"
)

func main() {
	logger.InitLogger()

	cfg := config.LoadConfig()
	connection := app.NewDB(cfg)

	userRepo := repository.NewUserRepository(connection.DB)
	userService := service.NewUserService(userRepo, cfg)

	noteRepo := repository.NewNoteRepository(connection.DB)
	noteService := service.NewNoteService(noteRepo)

	router := app.NewRouter(userService, noteService, cfg)
	http.ListenAndServe(":8080", router)
}
