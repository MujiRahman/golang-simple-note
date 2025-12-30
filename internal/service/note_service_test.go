package service

import (
	"errors"
	"testing"

	"github.com/MujiRahman/golang-simple-note/internal/model"
)

type mockNoteRepo struct {
	notes  map[uint]*model.Note
	nextID uint
}

func newMockNoteRepo() *mockNoteRepo {
	return &mockNoteRepo{notes: make(map[uint]*model.Note), nextID: 1}
}

func (m *mockNoteRepo) Create(note *model.Note) error {
	if note.ID == 0 {
		note.ID = m.nextID
		m.nextID++
	}
	m.notes[note.ID] = note
	return nil
}

func (m *mockNoteRepo) FindByID(id uint) (*model.Note, error) {
	n, ok := m.notes[id]
	if !ok {
		return nil, nil
	}
	return n, nil
}

func (m *mockNoteRepo) FindByUser(userID uint) ([]model.Note, error) {
	var out []model.Note
	for _, n := range m.notes {
		if n.UserID == userID {
			out = append(out, *n)
		}
	}
	return out, nil
}

func (m *mockNoteRepo) Update(note *model.Note) error {
	if _, ok := m.notes[note.ID]; !ok {
		return errors.New("not found")
	}
	m.notes[note.ID] = note
	return nil
}

func (m *mockNoteRepo) Delete(id uint) error {
	if _, ok := m.notes[id]; !ok {
		return errors.New("not found")
	}
	delete(m.notes, id)
	return nil
}

func TestNoteService_CRUD(t *testing.T) {
	repo := newMockNoteRepo()
	svc := NewNoteService(repo)

	// Create
	n, err := svc.Create(10, "t1", "c1")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if n.UserID != 10 || n.Title != "t1" {
		t.Fatalf("Create returned unexpected note: %+v", n)
	}

	// GetByID success
	got, err := svc.GetByID(10, n.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.ID != n.ID {
		t.Fatalf("GetByID returned wrong note id: %d", got.ID)
	}

	// GetByID access denied
	_, err = svc.GetByID(11, n.ID)
	if err == nil {
		t.Fatalf("expected error when accessing note with wrong user")
	}

	// ListByUser
	notes, err := svc.ListByUser(10)
	if err != nil {
		t.Fatalf("ListByUser failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("ListByUser returned unexpected length: %d", len(notes))
	}

	// Update
	updated, err := svc.Update(10, n.ID, "t2", "c2")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Title != "t2" {
		t.Fatalf("Update did not change title: %v", updated.Title)
	}

	// Delete
	if err := svc.Delete(10, n.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	// confirm deleted
	got2, _ := svc.GetByID(10, n.ID)
	if got2 != nil {
		t.Fatalf("expected note to be deleted")
	}
}
