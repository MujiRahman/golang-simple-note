package app

import (
	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/repository"
	"github.com/MujiRahman/golang-simple-note/internal/service"
)

// Container groups repositories and services for easy DI and maintenance.
type Repositories struct {
	User repository.UserRepository
	Note repository.NoteRepository
}

type Services struct {
	User service.UserService
	Note service.NoteService
}

type Container struct {
	Repos Repositories
	Svcs  Services
}

// NewContainer wires repositories and services using the provided DB connection and config.
func NewContainer(conn *Connect, cfg *config.Config) *Container {
	userRepo := repository.NewUserRepository(conn.DB)
	noteRepo := repository.NewNoteRepository(conn.DB)

	userSvc := service.NewUserService(userRepo, cfg)
	noteSvc := service.NewNoteService(noteRepo)

	return &Container{
		Repos: Repositories{User: userRepo, Note: noteRepo},
		Svcs:  Services{User: userSvc, Note: noteSvc},
	}
}
