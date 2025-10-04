package service

import (
	"errors"

	"github.com/MujiRahman/golang-simple-note/internal/model"
	"github.com/MujiRahman/golang-simple-note/internal/repository"
)

type NoteService interface {
	Create(userID uint, title, content string) (*model.Note, error)
	GetByID(userID, id uint) (*model.Note, error)
	ListByUser(userID uint) ([]model.Note, error)
	Update(userID, id uint, title, content string) (*model.Note, error)
	Delete(userID, id uint) error
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) Create(userID uint, title, content string) (*model.Note, error) {
	n := &model.Note{
		UserID:  userID,
		Title:   title,
		Content: content,
	}
	if err := s.repo.Create(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *noteService) GetByID(userID, id uint) (*model.Note, error) {
	n, err := s.repo.FindByID(id)
	if err != nil || n == nil {
		return nil, err
	}
	if n.UserID != userID {
		return nil, errors.New("not found or access denied")
	}
	return n, nil
}

func (s *noteService) ListByUser(userID uint) ([]model.Note, error) {
	return s.repo.FindByUser(userID)
}

func (s *noteService) Update(userID, id uint, title, content string) (*model.Note, error) {
	n, err := s.repo.FindByID(id)
	if err != nil || n == nil {
		return nil, err
	}
	if n.UserID != userID {
		return nil, errors.New("not found or access denied")
	}
	n.Title = title
	n.Content = content
	if err := s.repo.Update(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *noteService) Delete(userID, id uint) error {
	n, err := s.repo.FindByID(id)
	if err != nil || n == nil {
		return err
	}
	if n.UserID != userID {
		return errors.New("not found or access denied")
	}
	return s.repo.Delete(id)
}
