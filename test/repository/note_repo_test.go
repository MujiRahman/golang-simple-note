package repository_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MujiRahman/golang-simple-note/internal/model"
	"github.com/MujiRahman/golang-simple-note/internal/repository"
	"github.com/MujiRahman/golang-simple-note/test/testutil"
)

func TestNoteRepository_CreateFindByID(t *testing.T) {
	gdb, mock, sqlDB, err := testutil.NewGormWithSqlmock()
	if err != nil {
		t.Fatalf("failed create gorm+sqlmock: %v", err)
	}
	defer sqlDB.Close()

	repo := repository.NewNoteRepository(gdb)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `notes`")).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	n := &model.Note{UserID: 1, Title: "T", Content: "C"}
	if err := repo.Create(n); err != nil {
		t.Fatalf("create note failed: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at"}).AddRow(1, 1, "T", "C", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*notes.*WHERE .*LIMIT \\?").WillReturnRows(rows)

	got, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID error: %v", err)
	}
	if got == nil || got.ID != 1 {
		t.Fatalf("expected note id 1, got %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
