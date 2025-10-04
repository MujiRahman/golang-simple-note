package main

import (
	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/app"
	"github.com/MujiRahman/golang-simple-note/internal/controller"
	"github.com/MujiRahman/golang-simple-note/internal/repository"
	"github.com/MujiRahman/golang-simple-note/internal/service"
	"github.com/MujiRahman/golang-simple-note/pkg/logger"
)

func main() {
	logger.InitLogger()

	cfg := config.LoadConfig()
	connection := app.NewDB(cfg)

	noteRepo := repository.NewNoteRepository(connection.DB)
	noteService := service.NewNoteService(noteRepo)
	noteController := controller.NewNoteController(noteService)

	router := app.NewRouter(noteController)
	router.Run(":8080")
}
