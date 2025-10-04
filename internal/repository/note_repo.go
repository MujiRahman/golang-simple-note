package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/MujiRahman/golang-simple-note/internal/model"
)

type NoteRepository interface {
	Create(note *model.Note) error
	FindByID(id uint) (*model.Note, error)
	FindByUser(userID uint) ([]model.Note, error)
	Update(note *model.Note) error
	Delete(id uint) error
}

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(note *model.Note) error {
	return r.db.Create(note).Error
}

func (r *noteRepository) FindByID(id uint) (*model.Note, error) {
	var n model.Note
	if err := r.db.First(&n, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func (r *noteRepository) FindByUser(userID uint) ([]model.Note, error) {
	var notes []model.Note
	if err := r.db.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *noteRepository) Update(note *model.Note) error {
	return r.db.Save(note).Error
}

func (r *noteRepository) Delete(id uint) error {
	return r.db.Delete(&model.Note{}, id).Error
}
