package main

import (
	"net/http"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/app"
	"github.com/MujiRahman/golang-simple-note/pkg/logger"
)

func main() {
	logger.InitLogger()

	cfg := config.LoadConfig()
	connection := app.NewDB(cfg)

	// centralize wiring of repos & services
	container := app.NewContainer(connection, cfg)
	userService := container.Svcs.User
	noteService := container.Svcs.Note

	router := app.NewRouter(userService, noteService, cfg)
	http.ListenAndServe(":8080", router)
}
