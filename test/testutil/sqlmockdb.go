package testutil

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewGormWithSqlmock creates a *gorm.DB backed by sqlmock and returns the sqlmock
// instance and the underlying *sql.DB (so caller can Close it when done).
func NewGormWithSqlmock() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		db.Close()
		return nil, nil, nil, err
	}

	return gdb, mock, db, nil
}
