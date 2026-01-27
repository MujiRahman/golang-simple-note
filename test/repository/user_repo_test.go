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

func TestUserRepository_CreateFindByUsername(t *testing.T) {
	gdb, mock, sqlDB, err := testutil.NewGormWithSqlmock()
	if err != nil {
		t.Fatalf("failed create gorm+sqlmock: %v", err)
	}
	defer sqlDB.Close()

	repo := repository.NewUserRepository(gdb)

	// Expect insert (GORM starts transaction)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	u := &model.User{Name: "Test", Username: "tester", Password: "hashed", Email: "a@a.com"}
	if err := repo.Create(u); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Expect select by username
	rows := sqlmock.NewRows([]string{"id", "name", "username", "password", "email", "created_at", "updated_at"}).AddRow(1, "Test", "tester", "hashed", "a@a.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*users.*WHERE .*username.*LIMIT \\?").WillReturnRows(rows)

	got, err := repo.FindByUsername("tester")
	if err != nil {
		t.Fatalf("FindByUsername error: %v", err)
	}
	if got == nil || got.Username != "tester" {
		t.Fatalf("expected username tester, got %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
